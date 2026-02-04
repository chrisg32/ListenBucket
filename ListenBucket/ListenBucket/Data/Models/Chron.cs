using System.ComponentModel.DataAnnotations;

namespace ListenBucket.Data.Models;

public class Chron
{
    public int Id { get; set; }
    [MaxLength(3)]
    public string? Minute { get; set; }
    [MaxLength(3)]
    public string? Hour { get; set; }
    [MaxLength(3)]
    public string? Day { get; set; }
    [MaxLength(3)]
    public string? Month { get; set; }
    [MaxLength(3)]
    public string? DayOfWeek { get; set; }

    public override string ToString()
    {
        return $"{F(Minute)} {F(Hour)} {F(Day)} {F(Month)} {F(DayOfWeek)}";
    }
    private static string F(string? part)
    {
        return string.IsNullOrWhiteSpace(part) ? "*" : part;
    }
}