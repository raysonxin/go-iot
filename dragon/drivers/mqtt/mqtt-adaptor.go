package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	paho "github.com/eclipse/paho.mqtt.golang"
	multierror "github.com/hashicorp/go-multierror"
)

type Message paho.Message

// Adaptor is the mqtt client wrapper
type Adaptor struct {
	Name          string
	Host          string
	ClientID      string
	UserName      string
	Password      string
	UseSSL        bool
	ServerCert    string
	ClientCert    string
	ClientKey     string
	AutoReconnect bool
	CleanSession  bool
	Client        paho.Client
}

func NewAdaptor(host string, clientID string) *Adaptor {
	return &Adaptor{
		Name:          "mqtt",
		Host:          host,
		AutoReconnect: true,
		CleanSession:  true,
		UseSSL:        false,
		ClientID:      clientID,
	}
}

func NewAdaptorWithAuth(host, clientID, username, password string) *Adaptor {
	return &Adaptor{
		Name:          "mqtt",
		Host:          host,
		AutoReconnect: true,
		CleanSession:  true,
		UseSSL:        false,
		ClientID:      clientID,
		UserName:      username,
		Password:      password,
	}
}

func (a *Adaptor) SetCerts(serverCert, clientCert, clientKey string) {
	a.ServerCert = serverCert
	a.ClientCert = clientCert
	a.ClientKey = clientKey
}

func (a *Adaptor) Connect() (err error) {
	a.Client = paho.NewClient(a.createOptions())
	if token := a.Client.Connect(); token.Wait() && token.Error() != nil {
		err = multierror.Append(err, token.Error())
	}
	return
}

func (a *Adaptor) Disconnect() (err error) {
	if a.Client != nil {
		a.Client.Disconnect(500)
	}
	return
}

func (a *Adaptor) Finalize() (err error) {
	a.Disconnect()
	return
}

func (a *Adaptor) Publish(topic string, message []byte) bool {
	if a.Client == nil {
		return false
	}
	a.Client.Publish(topic, 0, false, message)
	return true
}

func (a *Adaptor) Subscribe(topic string, f func(msg Message)) bool {
	if a.Client == nil {
		return false
	}
	a.Client.Subscribe(topic, 0, func(client paho.Client, msg paho.Message) {
		f(msg)
	})
	return true
}

func (a *Adaptor) createOptions() *paho.ClientOptions {
	opts := paho.NewClientOptions()
	opts.AddBroker(a.Host) //(a.Host)
	opts.SetClientID(a.ClientID)

	if a.UserName != "" && a.Password != "" {
		opts.SetPassword(a.Password)
		opts.SetUsername(a.UserName)
	}

	opts.AutoReconnect = a.AutoReconnect
	opts.CleanSession = a.CleanSession

	if a.UseSSL {
		opts.SetTLSConfig(a.createTLSConfig())
	}

	return opts
}

func (a *Adaptor) createTLSConfig() *tls.Config {
	var certpool *x509.CertPool
	if len(a.ServerCert) > 0 {
		certpool = x509.NewCertPool()
		pemCerts, err := ioutil.ReadFile(a.ServerCert)
		if err == nil {
			certpool.AppendCertsFromPEM(pemCerts)
		}
	}

	var certs []tls.Certificate
	if len(a.ClientCert) > 0 && len(a.ClientKey) > 0 {
		cert, err := tls.LoadX509KeyPair(a.ClientCert, a.ClientKey)
		if err != nil {
			//TODO::
			panic(err)
		}
		certs = append(certs, cert)
	}

	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       certs,
	}
}
