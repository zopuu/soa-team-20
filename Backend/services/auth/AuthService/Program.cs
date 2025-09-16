using AuthService.Data;
using AuthService.Repositories;
using AuthService.Services;
using AuthService.Utills;
using Microsoft.AspNetCore.Authentication.JwtBearer;
using Microsoft.EntityFrameworkCore;
using Microsoft.IdentityModel.Tokens;
using Microsoft.OpenApi.Models;
using System.Text;
using System.Text.Json.Serialization;
using Serilog;
using Serilog.Events;
using Serilog.Formatting.Compact;

var builder = WebApplication.CreateBuilder(args);

builder.Logging.ClearProviders();

Log.Logger = new LoggerConfiguration()
    .MinimumLevel.Override("Microsoft", LogEventLevel.Warning)
    .MinimumLevel.Is(Enum.Parse<LogEventLevel>(
        Environment.GetEnvironmentVariable("LOG_LEVEL") ?? "Information", true))
    .Enrich.WithProperty("service", "auth")
    .Enrich.FromLogContext()
    .Enrich.WithEnvironmentName()
    .Enrich.WithMachineName()
    .WriteTo.Console(new CompactJsonFormatter())
    .CreateLogger();

builder.Host.UseSerilog();


// CORS
builder.Services.AddCors(opts =>
    opts.AddPolicy("AllowAngularDevClient", policy =>
        policy.WithOrigins("http://localhost:4200")
            .AllowAnyHeader()
            .AllowAnyMethod()
    ));



// Add services to the container.

builder.Services.AddDbContext<AppDbContext>(opts =>
    opts.UseNpgsql(builder.Configuration.GetConnectionString("DefaultConnection")));
builder.Services.AddScoped<IUserRepository, UserRepository>();
builder.Services.AddScoped<IAuthService, AuthService.Services.AuthService>();
builder.Services.AddScoped<IAdminService, AdminService>();
builder.Services.AddSingleton<JwtUtils>();

//Configure JWT bearer Auth
var key = Encoding.UTF8.GetBytes(builder.Configuration["Jwt:Key"]);
var issuer = builder.Configuration["Jwt:Issuer"];
var audience = builder.Configuration["Jwt:Audience"];

builder.Services.AddAuthentication(options =>
    {
        options.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
        options.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
    })
    .AddJwtBearer(options =>
    {
        options.RequireHttpsMetadata = false;
        options.SaveToken = true;
        options.TokenValidationParameters = new TokenValidationParameters
        {
            ValidateIssuer = true,
            ValidateAudience = true,
            ValidateLifetime = true,
            ValidateIssuerSigningKey = true,
            ValidIssuer = issuer,
            ValidAudience = audience,
            IssuerSigningKey = new SymmetricSecurityKey(key)
        };
    });
builder.Services.AddAuthorization();
builder.Services.AddControllers()
                .AddJsonOptions(o =>
                      o.JsonSerializerOptions.Converters.Add(new JsonStringEnumConverter()));
// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();
// 1) Add and configure Swagger with JWT support
builder.Services.AddSwaggerGen(c => {
    c.SwaggerDoc("v1", new OpenApiInfo { Title = "AuthService", Version = "v1" });

    // 1.a) Define the BearerAuth scheme
    c.AddSecurityDefinition("Bearer", new OpenApiSecurityScheme {
        Description = "JWT Bearer token — enter just the token (no 'Bearer ' prefix)",
        Name = "Authorization",
        In = ParameterLocation.Header,
        Type = SecuritySchemeType.Http,
        Scheme = "bearer",
        BearerFormat = "JWT"
    });

    // 1.b) Make all operations require the Bearer token by default
    c.AddSecurityRequirement(new OpenApiSecurityRequirement {
        {
            new OpenApiSecurityScheme {
                Reference = new OpenApiReference {
                    Type = ReferenceType.SecurityScheme,
                    Id   = "Bearer"
                }
            },
            Array.Empty<string>()
        }
    });
});

var app = builder.Build();
// apply any pending migrations
using (var scope = app.Services.CreateScope()) {
    var db = scope.ServiceProvider.GetRequiredService<AppDbContext>();
    var max = 10;
    var wait = TimeSpan.FromSeconds(5);

    for (int i = 0; i < max; i++) {
        try {
            db.Database.Migrate();
            break; // success!
        }
        catch (Npgsql.NpgsqlException) {
            if (i == max - 1) throw;     // re-throw after last try
            Thread.Sleep(wait);           // wait then retry
        }
    }
}

app.UseSerilogRequestLogging(opts => {
    opts.EnrichDiagnosticContext = (diag, http) => {
        diag.Set("request_id", http.Items["RequestId"]);
        diag.Set("trace_id", http.TraceIdentifier); // za sad fallback
        diag.Set("user_id", http.User?.Identity?.Name);
        diag.Set("method", http.Request.Method);
        diag.Set("path", http.Request.Path);
        diag.Set("client_ip", http.Connection.RemoteIpAddress?.ToString());
        diag.Set("scheme", http.Request.Scheme);
        diag.Set("host", http.Request.Host.ToString());
    };
});


//request_id middleware
app.Use(async (ctx, next) => {
    var reqId = ctx.Request.Headers["X-Request-ID"].FirstOrDefault()
            ?? Guid.NewGuid().ToString("N");
    ctx.Items["RequestId"] = reqId;
    ctx.Response.Headers["X-Request-ID"] = reqId;
    await next();
});

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI(c => {
        c.SwaggerEndpoint("/swagger/v1/swagger.json", "AuthService v1");
    });
}

app.UseRouting();
app.UseCors("AllowAngularDevClient");
app.UseAuthentication();
app.UseAuthorization();

app.MapControllers();
app.MapGet("/health", () => Results.Ok("Ok"));

app.Run();
