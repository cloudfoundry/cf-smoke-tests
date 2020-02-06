using System;

namespace NetSimple
{
    using Nancy;

    public class IndexModule : NancyModule
    {
        public IndexModule()
        {
            Get["/"] = x =>
            {

                var response = (Response)string.Format(
@"Healthy
It just needed to be restarted!
My application metadata: {0}
My port: {1}
My instance index: {2}
My custom env variable: {3}",
                    Environment.GetEnvironmentVariable("VCAP_APPLICATION"),
                    Environment.GetEnvironmentVariable("PORT"),
                    Environment.GetEnvironmentVariable("CF_INSTANCE_INDEX"),
                    Environment.GetEnvironmentVariable("CUSTOM_VAR")
                );
                response.ContentType = "text/plain";
                return response;
            };
        }
    }
}
