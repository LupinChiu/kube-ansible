using SockNet.ServerSocket;
using System;
using System.Collections.Generic;
using System.Net;
using System.Net.Sockets;
using System.Text;
using System.Threading.Tasks;

namespace EchoServer
{
    class Program
    {
        static void Main(string[] args)
        {
            var socketServer = new SocketServer();
            socketServer.InitializeSocketServer("0.0.0.0", int.Parse(args[0]));
            socketServer.SetReaderBufferBytes(1024);
            socketServer.StartListening();

            Console.WriteLine("Welcome Echo Server.");

            bool openServer = true;
            while (openServer)
            {
                if (socketServer.IsNewData())
                {
                    var data = socketServer.GetData();
                    // Do whatever you want with data
                    Task.Run(() => DoSomething(data, socketServer));
                }
            }

            //.... 
            socketServer.CloseServer();
        }
        private static void DoSomething(KeyValuePair<TcpClient, byte[]> data, SocketServer server)
        {
            Console.WriteLine(((IPEndPoint)data.Key.Client.RemoteEndPoint).Address.ToString() + ": " + Encoding.UTF8.GetString(data.Value));
            server.ResponseToClient(data.Key, "received");
        }
    }
}