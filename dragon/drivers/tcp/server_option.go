package tcp

type options struct {
	setCodecFunc func() MessageCodec
	onConnect    onConnectFunc
	onMessage    onMessageFunc
	onClose      onCloseFunc
	onError      onErrorFunc
	bufferSize   int
}

type ServerOption func(*options)

func OnBufferSizeOption(bsize int) ServerOption {
	return func(o *options) {
		o.bufferSize = bsize
	}
}

func OnConnectOption(cb func(Socket)) ServerOption {
	return func(o *options) {
		o.onConnect = cb
	}
}

func OnMessageOption(cb func(Message, Socket)) ServerOption {
	return func(o *options) {
		o.onMessage = cb
	}
}

func OnCloseOption(cb func(Socket)) ServerOption {
	return func(o *options) {
		o.onClose = cb
	}
}

func OnErrorOption(cb func(Socket)) ServerOption {
	return func(o *options) {
		o.onError = cb
	}
}

func SetCodecFuncOption(f func() MessageCodec) ServerOption {
	return func(o *options) {
		o.setCodecFunc = f
	}
}
