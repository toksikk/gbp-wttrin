package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var PluginName = "wttrin"
var PluginVersion = ""
var PluginBuilddate = ""

const (
	baseURL   = "https://wttr.in/"
	apiSuffix = "?format=j1"
)

var weatherCodes = map[string]string{
	"113": "☀️",   // Sunny
	"116": "⛅",    // Partly Cloudy
	"119": "☁️",   // Cloudy
	"122": "☁️",   // Very Cloudy
	"143": "🌫️",   // Fog
	"176": "🌦️",   // Light Showers
	"179": "🌨️",   // Light Sleet Showers
	"182": "🌨️",   // Light Sleet
	"185": "🌨️",   // Light Sleet
	"200": "⛈️",   // Thundery Showers
	"227": "🌨️",   // Light Snow
	"230": "❄️",   // Heavy Snow
	"248": "🌫️",   // Fog
	"260": "🌫️",   // Fog
	"263": "🌦️",   // Light Showers
	"266": "🌧️",   // Light Rain
	"281": "🌨️",   // Light Sleet
	"284": "🌨️",   // Light Sleet
	"293": "🌧️",   // Light Rain
	"296": "🌧️",   // Light Rain
	"299": "🌧️",   // Heavy Showers
	"302": "🌧️",   // Heavy Rain
	"305": "🌧️",   // Heavy Showers
	"308": "🌧️",   // Heavy Rain
	"311": "🌨️",   // Light Sleet
	"314": "🌨️",   // Light Sleet
	"317": "🌨️",   // Light Sleet
	"320": "🌨️",   // Light Snow
	"323": "🌨️",   // Light Snow Showers
	"326": "🌨️",   // Light Snow Showers
	"329": "❄️",   // Heavy Snow
	"332": "❄️",   // Heavy Snow
	"335": "❄️",   // Heavy Snow Showers
	"338": "❄️",   // Heavy Snow
	"350": "🌨️",   // Light Sleet
	"353": "🌦️",   // Light Showers
	"356": "🌧️",   // Heavy Showers
	"359": "🌧️",   // Heavy Rain
	"362": "🌨️",   // Light Sleet Showers
	"365": "🌨️",   // Light Sleet Showers
	"368": "🌨️",   // Light Snow Showers
	"371": "❄️",   // Heavy Snow Showers
	"374": "🌨️",   // Light Sleet Showers
	"377": "🌨️",   // Light Sleet
	"386": "⛈️",   // Thundery Showers
	"389": "⛈️",   // Thundery Heavy Rain
	"392": "❄️⛈️", // Thundery Snow Showers
	"395": "❄️",   // Heavy Snow Showers
}

var windDirectionEmojis = map[string]string{
	"N":   "⬆️",
	"NE":  "↗️",
	"E":   "➡️",
	"SE":  "↘️",
	"S":   "⬇️",
	"SW":  "↙️",
	"W":   "⬅️",
	"NW":  "↖️",
	"NNE": "⬆️",
	"ENE": "➡️",
	"ESE": "➡️",
	"SSE": "⬇️",
	"SSW": "⬇️",
	"WSW": "⬅️",
	"WNW": "⬅️",
	"NNW": "⬆️",
}

type wttrinResponse struct {
	CurrentCondition []struct {
		FeelsLikeC       string `json:"FeelsLikeC"`
		FeelsLikeF       string `json:"FeelsLikeF"`
		Cloudcover       string `json:"cloudcover"`
		Humidity         string `json:"humidity"`
		LocalObsDateTime string `json:"localObsDateTime"`
		ObservationTime  string `json:"observation_time"`
		PrecipInches     string `json:"precipInches"`
		PrecipMM         string `json:"precipMM"`
		Pressure         string `json:"pressure"`
		PressureInches   string `json:"pressureInches"`
		TempC            string `json:"temp_C"`
		TempF            string `json:"temp_F"`
		UvIndex          string `json:"uvIndex"`
		Visibility       string `json:"visibility"`
		VisibilityMiles  string `json:"visibilityMiles"`
		WeatherCode      string `json:"weatherCode"`
		WeatherDesc      []struct {
			Value string `json:"value"`
		} `json:"weatherDesc"`
		WeatherIconURL []struct {
			Value string `json:"value"`
		} `json:"weatherIconUrl"`
		Winddir16Point string `json:"winddir16Point"`
		WinddirDegree  string `json:"winddirDegree"`
		WindspeedKmph  string `json:"windspeedKmph"`
		WindspeedMiles string `json:"windspeedMiles"`
	} `json:"current_condition"`
	NearestArea []struct {
		AreaName []struct {
			Value string `json:"value"`
		} `json:"areaName"`
		Country []struct {
			Value string `json:"value"`
		} `json:"country"`
		Latitude   string `json:"latitude"`
		Longitude  string `json:"longitude"`
		Population string `json:"population"`
		Region     []struct {
			Value string `json:"value"`
		} `json:"region"`
		WeatherURL []struct {
			Value string `json:"value"`
		} `json:"weatherUrl"`
	} `json:"nearest_area"`
	Request []struct {
		Query string `json:"query"`
		Type  string `json:"type"`
	} `json:"request"`
	Weather []struct {
		Astronomy []struct {
			MoonIllumination string `json:"moon_illumination"`
			MoonPhase        string `json:"moon_phase"`
			Moonrise         string `json:"moonrise"`
			Moonset          string `json:"moonset"`
			Sunrise          string `json:"sunrise"`
			Sunset           string `json:"sunset"`
		} `json:"astronomy"`
		AvgtempC string `json:"avgtempC"`
		AvgtempF string `json:"avgtempF"`
		Date     string `json:"date"`
		Hourly   []struct {
			DewPointC        string `json:"DewPointC"`
			DewPointF        string `json:"DewPointF"`
			FeelsLikeC       string `json:"FeelsLikeC"`
			FeelsLikeF       string `json:"FeelsLikeF"`
			HeatIndexC       string `json:"HeatIndexC"`
			HeatIndexF       string `json:"HeatIndexF"`
			WindChillC       string `json:"WindChillC"`
			WindChillF       string `json:"WindChillF"`
			WindGustKmph     string `json:"WindGustKmph"`
			WindGustMiles    string `json:"WindGustMiles"`
			Chanceoffog      string `json:"chanceoffog"`
			Chanceoffrost    string `json:"chanceoffrost"`
			Chanceofhightemp string `json:"chanceofhightemp"`
			Chanceofovercast string `json:"chanceofovercast"`
			Chanceofrain     string `json:"chanceofrain"`
			Chanceofremdry   string `json:"chanceofremdry"`
			Chanceofsnow     string `json:"chanceofsnow"`
			Chanceofsunshine string `json:"chanceofsunshine"`
			Chanceofthunder  string `json:"chanceofthunder"`
			Chanceofwindy    string `json:"chanceofwindy"`
			Cloudcover       string `json:"cloudcover"`
			Humidity         string `json:"humidity"`
			PrecipInches     string `json:"precipInches"`
			PrecipMM         string `json:"precipMM"`
			Pressure         string `json:"pressure"`
			PressureInches   string `json:"pressureInches"`
			TempC            string `json:"tempC"`
			TempF            string `json:"tempF"`
			Time             string `json:"time"`
			UvIndex          string `json:"uvIndex"`
			Visibility       string `json:"visibility"`
			VisibilityMiles  string `json:"visibilityMiles"`
			WeatherCode      string `json:"weatherCode"`
			WeatherDesc      []struct {
				Value string `json:"value"`
			} `json:"weatherDesc"`
			WeatherIconURL []struct {
				Value string `json:"value"`
			} `json:"weatherIconUrl"`
			Winddir16Point string `json:"winddir16Point"`
			WinddirDegree  string `json:"winddirDegree"`
			WindspeedKmph  string `json:"windspeedKmph"`
			WindspeedMiles string `json:"windspeedMiles"`
		} `json:"hourly"`
		MaxtempC    string `json:"maxtempC"`
		MaxtempF    string `json:"maxtempF"`
		MintempC    string `json:"mintempC"`
		MintempF    string `json:"mintempF"`
		SunHour     string `json:"sunHour"`
		TotalSnowCm string `json:"totalSnow_cm"`
		UvIndex     string `json:"uvIndex"`
	} `json:"weather"`
}

func Start(discord *discordgo.Session) {
	discord.AddHandler(onMessageCreate)
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := strings.Replace(m.ContentWithMentionsReplaced(), s.State.Ready.User.Username, "username", 1)
	parts := strings.Split(msg, " ")

	if strings.ToLower(parts[0]) == "!wttr" {
		channel, err := s.State.Channel(m.ChannelID)
		if channel == nil {
			slog.Error("Failed to grab channel", "MessageID", m.ID, "ChannelID", m.ChannelID, "Error", err)
			return
		}

		guild, err := s.State.Guild(channel.GuildID)
		if guild == nil {
			slog.Error("Failed to grab guild", "MessageID", m.ID, "Channel", channel, "GuildID", channel.GuildID, "Error", err)
			return
		}

		handleWttrQuery(s, m, parts, guild)
	}
}

func handleWttrQuery(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild) {
	if len(parts) > 1 {
		location := strings.Join(parts[1:], "+")
		weatherResult, err := getWeather(location)
		if err != nil {
			slog.Error("Failed to get weather", "MessageID", m.ID, "Location", location, "Error", err)
			discordErrorMessage, err := s.ChannelMessageSend(m.ChannelID, "Failed to get weather for "+location+": "+err.Error())
			if err != nil {
				slog.Error("Failed to send message", "MessageID", discordErrorMessage.ID, "ChannelID", discordErrorMessage.ChannelID, "Error", err)
			}
			return
		}
		resultMessage := buildWeatherString(weatherResult)
		resultDiscordMessage, err := s.ChannelMessageSend(m.ChannelID, resultMessage)
		if err != nil {
			slog.Error("Failed to send message", "MessageID", resultDiscordMessage.ID, "ChannelID", resultDiscordMessage.ChannelID, "Error", err)
		}
	}
}

func buildWeatherString(weatherResult wttrinResponse) (result string) {
	weatherConditionEmoji := weatherCodes[weatherResult.CurrentCondition[0].WeatherCode]
	windDirectionEmoji := windDirectionEmojis[weatherResult.CurrentCondition[0].Winddir16Point]
	var region string
	if weatherResult.NearestArea[0].Region[0].Value != "" {
		region = "(" + weatherResult.NearestArea[0].Region[0].Value + ")"
	}

	r := "```📍 " + weatherResult.NearestArea[0].AreaName[0].Value + ", " + weatherResult.NearestArea[0].Country[0].Value + " " + region + "\n" +
		"🌡️ " + weatherResult.CurrentCondition[0].TempC + "°C (feels like " + weatherResult.CurrentCondition[0].FeelsLikeC + "°C)\n" +
		"💧 " + weatherResult.CurrentCondition[0].Humidity + "% humidity\n" +
		"🌬️ " + windDirectionEmoji + " " + weatherResult.CurrentCondition[0].WindspeedKmph + "km/h\n" +
		weatherConditionEmoji + " " + weatherResult.CurrentCondition[0].WeatherDesc[0].Value + "```"
	return r
}

func getWeather(location string) (weatherResult wttrinResponse, err error) {
	nocache := "&nonce=" + strconv.Itoa(rand.Intn(32768))
	queryURL := baseURL + location + apiSuffix + nocache
	slog.Info("Querying wttr.in", "URL", queryURL)
	return httpGet(queryURL)
}

func httpGet(url string) (weatherResult wttrinResponse, err error) {
	var resp *http.Response
	var httpClient = &http.Client{
		Timeout: 10 * 1000000000,
	}
	resp, err = httpClient.Get(url)
	if err != nil {
		slog.Error("Failed to get weather", "URL", url, "Error", err)
		return
	}
	defer resp.Body.Close()

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", "URL", url, "Error", err)
		return
	}

	err = json.Unmarshal(body, &weatherResult)
	slog.Info("Got weather", "URL", url, "Response", weatherResult)
	return
}
