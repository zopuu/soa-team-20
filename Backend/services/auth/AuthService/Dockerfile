# Stage 1: build
FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /src
# Install the EF CLI tool so we can run migrations at container startup
RUN dotnet tool install --global dotnet-ef --version 8.0.* \
  && export PATH="$PATH:/root/.dotnet/tools"
COPY ["AuthService.csproj", "./"]
RUN dotnet restore "./AuthService.csproj"
COPY . .
RUN dotnet publish "./AuthService.csproj" -c Release -o /app/publish

# Stage 2: runtime
FROM mcr.microsoft.com/dotnet/aspnet:8.0 AS runtime
WORKDIR /app
COPY --from=build /app/publish .
ENV ASPNETCORE_URLS=http://+:8080
EXPOSE 8080
ENTRYPOINT ["dotnet", "AuthService.dll"]
