package server

import (
	"net"

	"github.com/jingweno/ssh"
	log "github.com/sirupsen/logrus"
	gossh "golang.org/x/crypto/ssh"
)

func New(privates []gossh.Signer, socketDir string, logger log.FieldLogger) *Server {
	return &Server{
		privates:  privates,
		socketDir: socketDir,
		logger:    logger,
	}
}

type Server struct {
	privates  []gossh.Signer
	socketDir string
	logger    log.FieldLogger
}

func (s *Server) Serve(ln net.Listener) error {
	sh := newStreamlocalForwardHandler(s.socketDir, s.logger.WithField("handler", "streamlocalForwardHandler"))
	ph := newSSHProxyHandler(s.socketDir, s.logger.WithField("handler", "sshProxyHandler"))

	// convert []gossh.Signer to []ssh.Signer
	// TODO: use gossh only
	var signers []ssh.Signer
	for _, p := range s.privates {
		signers = append(signers, p)
	}

	server := ssh.Server{
		HostSigners: signers,
		Handler:     ph.Handler,
		ReversePortForwardingCallback: ssh.ReversePortForwardingCallback(func(ctx ssh.Context, host string, port uint32) (granted bool) {
			s.logger.WithFields(log.Fields{"tunnel-host": host, "tunnel-port": port}).Info("attempt to bind")
			return true
		}),
		RequestHandlers: map[string]ssh.RequestHandler{
			streamlocalForwardChannelType:      sh.Handler,
			cancelStreamlocalForwardChanneType: sh.Handler,
			pingRequestType:                    pingRequestHandler,
		},
	}

	return server.Serve(ln)
}

func pingRequestHandler(ctx ssh.Context, srv *ssh.Server, req *gossh.Request) (bool, []byte) {
	return true, nil
}
