using FluentValidation;
using ListenBucket.Data;
using ListenBucket.Data.Models;
using ListenBucket.Services;
using Microsoft.AspNetCore.Components.Forms;
using Microsoft.EntityFrameworkCore;

namespace ListenBucket.ViewModels;

public class SettingsViewModel(ApplicationDbContext context, JobService jobService)
{
    public string? GoogleApiKey { get; set; }

    public List<Chron> PollChronSettings { get; set; } = new();
    public async Task Load()
    {
        var settings = await context.Settings.FirstOrDefaultAsync();
        settings ??= new Settings();
        GoogleApiKey = settings.GoogleApiKey;
        PollChronSettings = await context.PollChronSettings.ToListAsync();
    }

    public async Task SaveSettings(EditContext arg)
    {
        if(!arg.Validate() || !arg.IsModified())
            return;
        var settings = await context.Settings.FirstOrDefaultAsync();
        if (settings == null)
        {
            settings = new Settings();
            context.Settings.Add(settings);
        }
        settings.GoogleApiKey = GoogleApiKey;

        // add or update
        foreach (var chron in PollChronSettings)
        {
            if (chron.Id == 0)
            {
                context.PollChronSettings.Add(chron);
            }
            else
            {
                context.PollChronSettings.Update(chron);
            }
        }
        
        // remove
        var toRemove = await context.PollChronSettings.Where(c => c.Id != -1 && !PollChronSettings.Select(p => p.Id).Contains(c.Id)).ToListAsync();
        context.PollChronSettings.RemoveRange(toRemove);
        
        await context.SaveChangesAsync();
        
        var updatedChronSettings = await context.PollChronSettings.ToListAsync();
        
        jobService.RemovePollingIntervals(toRemove);
        jobService.AddPollingIntervals(updatedChronSettings);
    }


    public void AddPollChron()
    {
        PollChronSettings.Add(new Chron());
    }
}



public class SettingsViewModelValidator : AbstractValidator<SettingsViewModel>
{
    public SettingsViewModelValidator()
    {
        RuleFor(x => x.GoogleApiKey)
            .NotEmpty()
            .WithMessage("Google API Key is required.")
            .Matches(@"^[A-Za-z0-9_\-]{39}$")
            .WithMessage("Invalid Google API Key format.");
        
        RuleForEach(x => x.PollChronSettings)
            .SetValidator(new ChronViewModelValidator());
    }
}

public class ChronViewModelValidator : AbstractValidator<Chron>
{
    public ChronViewModelValidator()
    {
        const string match = @"^(\d{1,2}|\*|)$";
        
        RuleFor(x => x.Minute)
            .Matches(match)
            .WithMessage("Invalid Minute format.");
        
        RuleFor(x => x.Hour)
            .Matches(match)
            .WithMessage("Invalid Hour format.");
        
        RuleFor(x => x.Day)
            .Matches(match)
            .WithMessage("Invalid Day format.");
        
        RuleFor(x => x.Month)
            .Matches(match)
            .WithMessage("Invalid Month format.");
        
        RuleFor(x => x.DayOfWeek)
            .Matches(match)
            .WithMessage("Invalid Day of Week format.");
    }
}