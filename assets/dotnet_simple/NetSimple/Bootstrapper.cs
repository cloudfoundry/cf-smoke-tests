using System;
using System.Diagnostics;
using System.Threading;
using Nancy;
using Nancy.Bootstrapper;
using Nancy.TinyIoc;

namespace NetSimple
{
    public class Bootstrapper : DefaultNancyBootstrapper
    {
        protected override void ApplicationStartup(TinyIoCContainer container, IPipelines pipelines)
        {
            new Thread(() =>
            {
                while (true)
                {
                    var unixTimestamp = (Int32)(DateTime.UtcNow.Subtract(new DateTime(1970, 1, 1))).TotalSeconds;
                    Console.WriteLine("Tick: {0}", unixTimestamp);
                    Thread.Sleep(TimeSpan.FromSeconds(1));
                }
            }).Start();
        }
    }
}
