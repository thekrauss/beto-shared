package eventbus

func (r *RabbitMQ) DeclareExchange(name, kind string) error {
	return r.channel.ExchangeDeclare(
		name,
		kind,  // "direct", "fanout", "topic", "headers"
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // args
	)
}

func (r *RabbitMQ) DeclareQueue(name string) error {
	_, err := r.channel.QueueDeclare(
		name,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	return err
}

func (r *RabbitMQ) BindQueue(queue, routingKey, exchange string) error {
	return r.channel.QueueBind(
		queue,
		routingKey,
		exchange,
		false,
		nil,
	)
}
