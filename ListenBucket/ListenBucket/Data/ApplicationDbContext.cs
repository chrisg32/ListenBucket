using ListenBucket.Data.Models;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;

namespace ListenBucket.Data;

public class ApplicationDbContext(DbContextOptions<ApplicationDbContext> options)
    : IdentityDbContext<ApplicationUser>(options)
{
    
    public DbSet<Podcast> Podcasts { get; set; }
    public DbSet<Item> Items { get; set; }
    public DbSet<Settings> Settings { get; set; }
    public DbSet<Chron> PollChronSettings { get; set; }

    public void Seed()
    {
        if (!Podcasts.Any())
        {
            Podcasts.Add(new Podcast
            {
                Id = 1,
                Title = "Listen Later",
                Description = "Podcast for items you want to listen to later",
                Source = null,
                YouTubeId = null,
                ThumbnailUrl = "/images/listen-later.png",
            });
            SaveChanges();
        }
    }
}