using System.ComponentModel.DataAnnotations;

namespace ListenBucket.Data.Models;

public class Settings
{
    public int Id { get; set; }
    
    [MaxLength(39)]
    public string? GoogleApiKey { get; set; }
}