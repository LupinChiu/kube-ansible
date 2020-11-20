using SockNet.ClientSocket;
using System;
using System.Net;
using System.Text;
using System.Threading;

namespace EchoClientDonet
{
    class Program
    {
        static async System.Threading.Tasks.Task Main(string[] args)
        {
            byte[] recData = null;
            SocketClient client = new SocketClient(args[0], int.Parse(args[1]));
            try
            {
                if (await client.Connect())
                {
                    Console.WriteLine("Connect to : {0}:{1}",args[0],args[1]);
                    var host = Dns.GetHostEntry(Dns.GetHostName());
                    while (true)
                    {
                        await client.Send($"{host.HostName} Send: {DateTime.Now.ToString()}");
                        recData = await client.ReceiveBytes();
                        Console.WriteLine("Received data: " + Encoding.UTF8.GetString(recData));
                        Thread.Sleep(int.Parse(args[2]));
                    }
                }
            }
            catch (Exception e)
            {
                Console.WriteLine("Exception raised: " + e);
            }
            //...
            client.Disconnect();
        }
    }
}
