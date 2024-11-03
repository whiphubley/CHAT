#!/usr/bin/ruby

require 'socket'
include Socket::Constants

def bind_socket(port)
  puts "the port is #{port}"
  server = Socket.new(AF_INET, SOCK_STREAM, 0)
  addr = Socket.sockaddr_in(port, '127.0.0.1')
  server.bind(addr)
  server.listen(5)
  puts "listening on port #{port}..."
  get_connection(server)
end

def get_connection(connection)
  loop do
    client, client_addr = connection.accept
    Thread.new(client, client_addr) do |conn, conn_addr|
      puts "client connected: #{conn_addr.inspect}"
      data = conn.recv(1024)
      puts "received data: #{data}"
      conn.puts "hello from rubychat\nyou said: #{data}"
      conn.close
    end
  end
end

def main
  bind_socket(3333)
end

main
