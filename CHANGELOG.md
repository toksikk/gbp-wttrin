## v1.4.0 (2024-03-22)

### Feat

- **formatting**: make location link to google maps and reduce redundancy
- **output**: add markdown for readability
- **warnings**: remove parantheses

## v1.3.0 (2023-12-06)

### Feat

- **warnings**: add warnings of high chances to wttr and wttrf
- **forecast**: improve condition output by using most occuring forecast weather codes
- **wttrin**: calculate averages for wind direction and wind speed in forecast
- **wttrin**: Add average rainfall calculation to buildForecastString function

### Fix

- **forecast**: line break when no rain but snow

### Refactor

- **wttrin**: remove stupid code / strconv.Atoi wrapper
- **wttrin**: windDirDegreeStringToInt function to convertStringToInt
- **wttrin**: wind direction conversion to use integer values
- **wttrin**: handleWttrQuery function to use boolean flag for forecast

## v1.2.1 (2023-12-06)

### Feat

- **wttrin**: add wind direction to forecast

## v1.2.0 (2023-12-06)

### Feat

- **wttrin**: add !wttrf for a forecast

### Refactor

- **wttrin**: enhance forecast output for readability
- **wttrin**: extract weatherConditionEmoji definition to separate function
- **wttrin**: extract windir emoji definition to separate function

## v1.1.0 (2023-11-30)

### Feat

- **lib**: move logic to importable plugin package
- **logging**: add logging call for unknown weather codes
- **wttrin**: improve things
- **output/log**: nicer output and add some info logging
- **rewrite**: complete rewrite to use JSON api of wttr.in
- **build**: add update dependencies option to Makefile
- **wttrin**: avoid discord replacing unicode chars with emojis

### Fix

- **wttrin**: fix status code check
- **wttrin**: return error if requested location was not found
- **wttrin**: north emoji
- **service**: circumvent caching for more accurate data
- **build**: add missing ldflags to make release

## v0.1.1 (2022-04-13)

### Fix

- spaces in city names
