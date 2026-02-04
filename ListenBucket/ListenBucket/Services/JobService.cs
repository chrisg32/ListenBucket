using FFMpegCore;
using Hangfire;
using ListenBucket.Data;
using ListenBucket.Data.Models;
using Microsoft.EntityFrameworkCore;

namespace ListenBucket.Services;

public class JobService(IBackgroundJobClient backgroundJobClient, ApplicationDbContext context)
{
    public void AddPollingIntervals(List<Chron> chronSettings)
    {
        foreach (var chronSetting in chronSettings)
        {
            RecurringJob.AddOrUpdate($"polling-{chronSetting.Id}", () => Job_Poll(), chronSetting.ToString());
        }
    }
    
    public void RemovePollingIntervals(List<Chron> chronSettings)
    {
        foreach (var chronSetting in chronSettings)
        {
            RecurringJob.RemoveIfExists($"polling-{chronSetting.Id}");
        }
    }

    public async Task Job_Poll()
    {
        Console.WriteLine("============================================");
        Console.WriteLine("=============== I BE POLLING ===============");
        Console.WriteLine("============================================");
    }

    public string QueueDownload(int itemId)
    {
        return backgroundJobClient.Enqueue(() => Job_DownloadItem(itemId));
    }
    
    // ReSharper disable once MemberCanBePrivate.Global
    public async Task Job_DownloadItem(int itemId)
    {
        var item = await context.Items.FirstAsync(i => i.Id == itemId);
        var directory = Directory.GetCurrentDirectory();
        
        //download the file
        await DownloadService.DownloadAsync(item.Source, directory, item.YouTubeId).ContinueWith(async task =>
        {
            try
            {
                if (task.IsFaulted)
                {
                    // Handle error
                    Console.WriteLine($"Error downloading {item.Source}: {task.Exception}");
                    throw task.Exception;
                }

                var path = task.Result;

                if (string.IsNullOrEmpty(path))
                {
                    // Handle error
                    Console.WriteLine($"Error downloading {item.Source}: No path returned");
                    return;
                }

                //get the duration of the file
                var mediaInfo = await FFProbe.AnalyseAsync(path);

                //update the item with the path and duration
                item.Downloaded = true;
                item.Duration = mediaInfo.Duration;

                //update the item in the database
                context.Items.Update(item);
                await context.SaveChangesAsync();
                
            }catch (Exception e)
            {
                Console.WriteLine(e);
                throw;
            }
        });
    }
}