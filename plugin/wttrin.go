package gbpwttrin

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// PluginName is the name of the plugin
var PluginName = "wttrin"

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

// Start the plugin
func Start(discord *discordgo.Session) {
	discord.AddHandler(onMessageCreate)
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := strings.Replace(m.ContentWithMentionsReplaced(), s.State.Ready.User.Username, "username", 1)
	parts := strings.Split(msg, " ")
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

	switch strings.ToLower(parts[0]) {
	case "!wttr":
		handleWttrQuery(s, m, parts, guild, false)
	case "!wttrf":
		handleWttrQuery(s, m, parts, guild, true)
	}
}

func handleWttrQuery(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, g *discordgo.Guild, forecast bool) {
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

		var resultMessage string
		if forecast {
			resultMessage = buildForecastString(weatherResult)
		} else {
			resultMessage = buildWeatherString(weatherResult)
		}

		resultDiscordMessage, err := s.ChannelMessageSend(m.ChannelID, resultMessage)
		if err != nil {
			slog.Error("Failed to send message", "MessageID", resultDiscordMessage.ID, "ChannelID", resultDiscordMessage.ChannelID, "Error", err)
		}
	}
}

func getWindDirectionEmoji(winddirDegree int) (windDirectionEmoji string) {
	if (winddirDegree >= 337 && winddirDegree <= 360) || (winddirDegree >= 0 && winddirDegree <= 22) {
		windDirectionEmoji = "⬆️"
	} else if winddirDegree >= 22 && winddirDegree <= 67 {
		windDirectionEmoji = "↗️"
	} else if winddirDegree >= 67 && winddirDegree <= 112 {
		windDirectionEmoji = "➡️"
	} else if winddirDegree >= 112 && winddirDegree <= 157 {
		windDirectionEmoji = "↘️"
	} else if winddirDegree >= 157 && winddirDegree <= 202 {
		windDirectionEmoji = "⬇️"
	} else if winddirDegree >= 202 && winddirDegree <= 247 {
		windDirectionEmoji = "↙️"
	} else if winddirDegree >= 247 && winddirDegree <= 292 {
		windDirectionEmoji = "⬅️"
	} else if winddirDegree >= 292 && winddirDegree <= 337 {
		windDirectionEmoji = "↖️"
	}
	return
}

func getWeatherConditionEmoji(weatherCode string) (weatherConditionEmoji string) {
	weatherConditionEmoji = "🌈"
	for code := range weatherCodes {
		if weatherCode == code {
			weatherConditionEmoji = weatherCodes[code]
			break
		}
	}

	if weatherConditionEmoji == "🌈" {
		slog.Warn("Unknown weather code", "Code", weatherCode)
	}
	return
}

func buildWeatherString(weatherResult wttrinResponse) (result string) {
	weatherConditionEmoji := getWeatherConditionEmoji(weatherResult.CurrentCondition[0].WeatherCode)
	windDirDegree, err := strconv.Atoi(weatherResult.CurrentCondition[0].WinddirDegree)
	if err != nil {
		slog.Error("Failed to convert wind direction to integer", "weatherResult.CurrentCondition[0].WinddirDegree", weatherResult.CurrentCondition[0].WinddirDegree, "Error", err)
		return
	}
	windDirectionEmoji := getWindDirectionEmoji(windDirDegree)

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

func buildForecastString(weatherResult wttrinResponse) (result string) {
	result += "```"
	for i, day := range weatherResult.Weather {
		if i > 0 {
			result += "---\n"
		}
		weatherConditionEmoji := getWeatherConditionEmoji(day.Hourly[0].WeatherCode)
		// TODO: calculate average wind direction for all hours in the day
		windDirDegree, err := strconv.Atoi(weatherResult.CurrentCondition[0].WinddirDegree)
		if err != nil {
			slog.Error("Failed to convert wind direction to integer", "weatherResult.CurrentCondition[0].WinddirDegree", weatherResult.CurrentCondition[0].WinddirDegree, "Error", err)
			return
		}
		windDirectionEmoji := getWindDirectionEmoji(windDirDegree)
		result += "📅 " + day.Date + "\n" +
			"🌡️ " + day.MaxtempC + "°C / " + day.MintempC + "°C\n" +
			"🌬️ " + windDirectionEmoji + " " + day.Hourly[0].WindspeedKmph + "km/h\n" +
			weatherConditionEmoji + " " + day.Hourly[0].WeatherDesc[0].Value + "\n"

		totalSnow, err := strconv.ParseFloat(day.TotalSnowCm, 32)
		if err != nil {
			slog.Warn("Failed to parse total snow", "TotalSnowCm", day.TotalSnowCm, "Error", err)
		} else {
			if totalSnow > 0.0 {
				result += "❄️ " + day.TotalSnowCm + "cm"
			}
		}

		totalRain := 0.0
		for _, hour := range day.Hourly {
			rain, err := strconv.ParseFloat(hour.PrecipMM, 32)
			if err != nil {
				slog.Warn("Failed to parse total rain", "PrecipMM", hour.PrecipMM, "Error", err)
			}
			totalRain += rain
		}

		if totalRain > 0.0 {
			averageRain := totalRain / float64(len(day.Hourly))
			if averageRain > 0.0 {
				if totalSnow > 0.0 {
					result += " / "
				}
				result += fmt.Sprintf("🌧️ %.2fmm\n", averageRain)
			}
		} else {
			result += "\n" // Add newline if no rain but snow
		}
	}
	result += "```"
	return
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
		Timeout: 30 * time.Second,
	}
	resp, err = httpClient.Get(url)
	if err != nil {
		slog.Error("Failed to get weather", "URL", url, "Error", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		slog.Info("Could not find requested location", "URL", url, "StatusCode", resp.StatusCode)
		err = fmt.Errorf(resp.Status)
		return
	}

	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body", "URL", url, "Error", err)
		return
	}

	err = json.Unmarshal(body, &weatherResult)
	slog.Debug("Got weather", "URL", url, "Response", weatherResult)
	return
}
