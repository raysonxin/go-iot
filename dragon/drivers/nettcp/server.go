package nettcp

type options struct {
	codec     MessageCodec
	onConnect onConnectFunc
	onMessage onMessageFunc
	onClose   onCloseFunc
	onError   onErrorFunc
}
