var builder = DistributedApplication.CreateBuilder(args);

builder.AddProject<Projects.ListenBucket>("listenbucket");

builder.Build().Run();
