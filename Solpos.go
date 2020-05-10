package solpos

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
	"time"
)

/*
*    Contains:
*        S_solpos     (computes solar position and intensity
*                      from time and place)
*
*            INPUTS:     (via posdata struct) year, daynum, hour,
*                        minute, second, latitude, longitude, timezone,
*                        intervl
*            OPTIONAL:   (via posdata struct) month, day, press, temp, tilt,
*                        aspect, function
*            OUTPUTS:    EVERY variable in the struct posdata
*                            (defined in solpos.h)
*
*                       NOTE: Certain conditions exist during which some of
*                       the output variables are undefined or cannot be
*                       calculated.  In these cases, the variables are
*                       returned with flag values indicating such.  In other
*                       cases, the variables may return a realistic, though
*                       invalid, value. These variables and the flag values
*                       or invalid conditions are listed below:
*
*                       amass     -1.0 at zenetr angles greater than 93.0
*                                 degrees
*                       ampress   -1.0 at zenetr angles greater than 93.0
*                                 degrees
*                       azim      invalid at zenetr angle 0.0 or latitude
*                                 +/-90.0 or at night
*                       elevetr   limited to -9 degrees at night
*                       etr       0.0 at night
*                       etrn      0.0 at night
*                       etrtilt   0.0 when cosinc is less than 0
*                       prime     invalid at zenetr angles greater than 93.0
*                                 degrees
*                       sretr     +/- 2999.0 during periods of 24 hour sunup or
*                                 sundown
*                       ssetr     +/- 2999.0 during periods of 24 hour sunup or
*                                 sundown
*                       ssha      invalid at the North and South Poles
*                       unprime   invalid at zenetr angles greater than 93.0
*                                 degrees
*                       zenetr    limited to 99.0 degrees at night
*
*        S_init       (optional initialization for all input parameters in
*                      the posdata struct)
*           INPUTS:     struct posdata*
*           OUTPUTS:    struct posdata*
*
*                     (Note: initializes the required S_solpos INPUTS above
*                      to out-of-bounds conditions, forcing the user to
*                      supply the parameters; initializes the OPTIONAL
*                      S_solpos inputs above to nominal values.)
*
 *
*    Martin Rymes
*    National Renewable Energy Laboratory
*    25 March 1998
*
*    27 April 1999 REVISION:  Corrected leap year in S_date.
*    13 January 2000 REVISION:  SMW converted to structure posdata parameter
*                               and subdivided into functions.
*    01 February 2001 REVISION: SMW corrected ecobli calculation
*                               (changed sign). Error is small (max 0.015 deg
*                               in calculation of declination angle)
*/

/* Solpos interface: Each comment begins with a 1-column letter code:
   I:  INPUT variable
   O:  OUTPUT variable
   T:  TRANSITIONAL variable used in the algorithm,
       of interest only to the solar radiation
       modelers, and available to you because you
       may be one of them.
   The FUNCTION column indicates which sub-function
   within solpos must be switched on using the
   "function" parameter to calculate the desired
   output variable.  All function codes are
   defined in the solpos.h file.  The default
   S_ALL switch calculates all output variables.
   Multiple functions may be or'd to create a
   composite function switch.  For example,
   (S_TST | S_SBCF). Specifying only the functions
   for required output variables may allow solpos
   to execute more quickly.
   The S_DOY mask works as a toggle between the
   input date represented as a day number (daynum)
   or as month and day.  To set the switch (to
   use daynum input), the function is or'd; to
   clear the switch (to use month and day input),
   the function is inverted and and'd.
   For example:
       pdat->function |= S_DOY (sets daynum input)
       pdat->function &= ~S_DOY (sets month and day input)
   Whichever date form is used, S_solpos will
   calculate and return the variables(s) of the
   other form.  See the soltest.c program for
   other examples. */
type Solpos interface {
	// Methods
	Calculate() error
	// helper function to get sunrise
	GetSunrise() time.Time
	// helper function to get sunset
	GetSunset() time.Time
	// using go builtin time functions
	Getdate() time.Time
	SetDate(dt time.Time)

	/* I/O: S_DOY Day of month (May 27 = 27, etc.) solpos will CALCULATE this by default,
	   or will optionally require it as input depending on the setting of the S_DOY function switch. */
	GetDay() int
	SetDay(day int)

	/* I/O: S_DOY Day number (day of year; Feb 1 = 32 ) solpos REQUIRES this by default, but
	   will optionally calculate it from month and day depending on the setting of the S_DOY function switch. */
	GetDaynum() int
	SetDaynum(daydaynum int)
	/* I: Switch to choose functions for desired output. */
	GetFunction() SPFunctions
	SetFunction(function SPFunctions)

	/* I: Hour of day, 0 - 23, DEFAULT = 12 */
	GetHour() int
	SetHour(hour int)
	/* I: Interval of a measurement period in seconds.  Forces solpos to use the time and date from the interval midpoint. The INPUT time (hour,
	   minute, and second) is assumed to be the END of the measurement interval. */
	GetInterval() int
	SetInterval(interval int)
	/* I: Minute of hour, 0 - 59, DEFAULT = 0 */
	GetMinute() int
	SetMinute(minute int)
	/* I/O: S_DOY  Month number (Jan = 1, Feb = 2, etc.) solpos will CALCULATE this by default,
	   or will optionally require it as input depending on the setting of the S_DOY function switch. */
	GetMonth() int
	SetMonth(month int)
	/* I: Second of minute, 0 - 59, DEFAULT = 0 */
	GetSecond() int
	SetSecond(second int)
	/* I:  4-digit year (2-digit year is NOT allowed */
	GetYear() int
	SetYear(year int)
	/* O:  S_AMASS    Relative optical airmass */
	GetAmass() float64
	/* O:  S_AMASS    Pressure-corrected airmass */
	GetAmpress() float64
	/* I: Azimuth of panel surface (direction it faces) N=0, E=90, S=180, W=270, DEFAULT = 180 */
	GetAspect() float64
	SetAspect(aspect float64)
	/* O:  S_SOLAZM   Solar azimuth angle:  N=0, E=90, S=180, W=270 */
	GetAzim() float64
	/* O:  S_TILT     Cosine of solar incidence angle on panel */
	GetCosinc() float64
	/* O:  S_REFRAC   Cosine of refraction corrected solar zenith angle */
	GetCoszen() float64
	/* T:  S_GEOM     Day angle (daynum*360/year-length) degrees */
	GetDayang() float64
	/* T:  S_GEOM     Declination--zenith angle of solar noon at equator, degrees NORTH */
	GetDeclin() float64
	/* T:  S_GEOM     Ecliptic longitude, degrees */
	GetEclong() float64
	/* T:  S_GEOM     Obliquity of ecliptic */
	GetEcobli() float64
	/* T:  S_GEOM     Time of ecliptic calculations */
	GetEctime() float64
	/* O:  S_ZENETR   Solar elevation, no atmospheric correction (= ETR) */
	GetElevetr() float64
	/* O:  S_REFRAC   Solar elevation angle, deg. from horizon, refracted */
	GetElevref() float64
	/* T:  S_TST      Equation of time (TST - LMT), minutes */
	GetEqntim() float64
	/* T:  S_GEOM     Earth radius vector (multiplied to solar constant) */
	GetErv() float64
	/* O:  S_ETR      Extraterrestrial (top-of-atmosphere) W/sq m global horizontal solar irradiance */
	GetEtr() float64
	/* O:  S_ETR      Extraterrestrial (top-of-atmosphere) W/sq m direct normal solar irradiance */
	GetEtrn() float64
	/* O:  S_TILT     Extraterrestrial (top-of-atmosphere) W/sq m global irradiance on a tilted surface */
	GetEtrtilt() float64
	/* T:  S_GEOM     Greenwich mean sidereal time, hours */
	GetGmst() float64
	/* T:  S_GEOM     Hour angle--hour of sun from solar noon, degrees WEST */
	GetHrang() float64
	/* T:  S_GEOM     Julian Day of 1 JAN 2000 minus 2,400,000 days (in order to regain single precision) */
	GetJulday() float64
	/* I: Latitude, degrees north (south negative) */
	GetLatitude() float64
	SetLatitude(latitude float64)
	/* I: Longitude, degrees east (west negative) */
	GetLongitude() float64
	SetLongitude(longitude float64)
	/* T:  S_GEOM     Local mean sidereal time, degrees */
	GetLmst() float64
	/* T:  S_GEOM     Mean anomaly, degrees */
	GetMnanom() float64
	/* T:  S_GEOM     Mean longitude, degrees */
	GetMnlong() float64
	/* T:  S_GEOM     Right ascension, degrees */
	GetRascen() float64
	/* I:             Surface pressure, millibars, used for refraction correction and ampress */
	GetPress() float64
	SetPress(press float64)
	/* O:  S_PRIME    Factor that normalizes Kt, Kn, etc. */
	GetPrime() float64
	/* O:  S_SBCF     Shadow-band correction factor */
	GetSbcf() float64
	/* I:             Shadow-band width (cm) */
	GetSbwid() float64
	SetSbwid(sbwid float64)
	/* I:             Shadow-band radius (cm) */
	GetSbrad() float64
	SetSbrad(sbrad float64)
	/* I:             Shadow-band sky factor */
	GetSbsky() float64
	SetSbsky(sbsky float64)
	/* I:             Solar constant (NREL uses 1367 W/sq m) */
	GetSolcon() float64
	SetSolcon(solcon float64)
	/* T:  S_SRHA     Sunset(/rise) hour angle, degrees */
	GetSsha() float64
	/* O:  S_SRSS     Sunrise time, minutes from midnight, local, WITHOUT refraction */
	GetSretr() float64
	/* O:  S_SRSS     Sunset time, minutes from midnight, local, WITHOUT refraction */
	GetSsetr() float64
	/* I:             Ambient dry-bulb temperature, degrees C, used for refraction correction */
	GetTemp() float64
	SetTemp(temp float64)
	/* I:             Degrees tilt from horizontal of panel */
	GetTilt() float64
	SetTilt(tilt float64)
	/* I:             Time zone, east (west negative). USA:  Mountain = -7, Central = -6, etc. */
	GetTimezone() float64
	SetTimezone(timezone float64)
	/* T:  S_TST      True solar time, minutes from midnight */
	GetTst() float64
	/* T:  S_TST      True solar time - local standard time */
	GetTstfix() float64
	/* O:  S_PRIME    Factor that denormalizes Kt', Kn', etc. */
	GetUnprime() float64
	/* T:  S_GEOM     Universal (Greenwich) standard time */
	GetUtime() float64
	/* T:  S_ZENETR   Solar zenith angle, no atmospheric correction (= ETR) */
	GetZenetr() float64
	/* O:  S_REFRAC   Solar zenith angle, deg. from zenith, refracted */
	GetZenref() float64
	SetZenref(zenref float64)
}

func NewSolpos(dt time.Time, latitude float64, longitude float64, optionalParameters map[string]interface{}) (Solpos, error) {
	var sp solpos
	sp.setTrigdata(trigdata{1.0, 1.0, 1.0, -999.0, 1.0})
	sp.init()
	sp.Latitude = latitude
	sp.Longitude = longitude
	sp.SetDate(dt)
	for key, value := range optionalParameters {
		switch key {
		case "press":
			tmpValue, ok := value.(float64)
			if !ok {
				err := errors.New("wrong type press, expected float64")
				return nil, err
			}
			sp.Press = tmpValue
		case "temp":
			tmpValue, ok := value.(float64)
			if !ok {
				err := errors.New("wrong type temp, expected float64")
				return nil, err
			}
			sp.Temp = tmpValue
		case "tilt":
			tmpValue, ok := value.(float64)
			if !ok {
				err := errors.New("wrong type tilt, expected float64")
				return nil, err
			}
			sp.Tilt = tmpValue

		case "aspect":
			tmpValue, ok := value.(float64)
			if !ok {
				err := errors.New("wrong type aspect, expected float64")
				return nil, err
			}
			sp.Aspect = tmpValue
		case "month":
			tmpValue, ok := value.(int)
			if !ok {
				err := errors.New("wrong type month, expected int")
				return nil, err
			}
			sp.Month = tmpValue
		case "day":
			tmpValue, ok := value.(int)
			if !ok {
				err := errors.New("wrong type day, expected int")
				return nil, err
			}
			sp.Day = tmpValue
		case "function":
			tmpValue, ok := value.(SPFunctions)
			if !ok {
				err := errors.New("wrong type, expected uint32")
				return nil, err
			}
			sp.Function = tmpValue
		}
	}
	return &sp, sp.Calculate()

}

type solpos struct {
	Day       int         // Day of month (May 27 = 27, etc.) solpos will CALCULATE this by default, or will optionally require it as input depending on the setting of the S_DOY  function switch.
	Daynum    int         // Day number (day of year; Feb 1 = 32 )	solpos REQUIRES this by default, but will optionally calculate it from month and day depending on the setting of the S_DOY function switch.
	Function  SPFunctions // Switch to choose functions for desired output.
	Hour      int         // Hour of day, 0 - 23, DEFAULT = 12 */
	Interval  int         // Interval of a measurement period in	seconds.  Forces solpos to use the time and date from the interval midpoint. The INPUT time (hour, minute, and second) is assumed to be the END of the measurement interval.
	Minute    int         // Minute of hour, 0 - 59, DEFAULT = 0
	Month     int         // Month number (Jan = 1, Feb = 2, etc.) solpos will CALCULATE this by default, or will optionally require it as input depending on the setting of the S_DOY function switch.
	Second    int         // Second of minute, 0 - 59, DEFAULT = 0
	Year      int         //  4-digit year (2-digit year is NOT  allowed
	Amass     float64     // Relative optical airmass */
	Ampress   float64     // Pressure-corrected airmass */
	Aspect    float64     // Azimuth of panel surface (direction it faces) N=0, E=90, S=180, W=270, DEFAULT = 180 */
	Azim      float64     // Solar azimuth angle:  N=0, E=90, S=180, W=270 */
	Cosinc    float64     // Cosine of solar incidence angle on panel */
	Coszen    float64     // Cosine of refraction corrected solar zenith angle */
	Dayang    float64     // Day angle (daynum*360/year-length) degrees */
	Declin    float64     // Declination--zenith angle of solar noon at equator, degrees NORTH */
	Eclong    float64     // Ecliptic longitude, degrees */
	Ecobli    float64     // Obliquity of ecliptic */
	Ectime    float64     // Time of ecliptic calculations */
	Elevetr   float64     // Solar elevation, no atmospheric correction (= ETR) */
	Elevref   float64     // Solar elevation angle, deg. from horizon, refracted */
	Eqntim    float64     // Equation of time (TST - LMT), minutes */
	Erv       float64     // Earth radius vector (multiplied to solar constant) */
	Etr       float64     // Extraterrestrial (top-of-atmosphere) W/sq m global horizontal solar irradiance */
	Etrn      float64     // Extraterrestrial (top-of-atmosphere) W/sq m direct normal solar irradiance */
	Etrtilt   float64     // Extraterrestrial (top-of-atmosphere) W/sq m global irradiance on a tilted surface */
	Gmst      float64     // Greenwich mean sidereal time, hours */
	Hrang     float64     // Hour angle--hour of sun from solar noon, degrees WEST */
	Julday    float64     // Julian Day of 1 JAN 2000 minus 2,400,000 days (in order to regain single precision) */
	Latitude  float64     // Latitude, degrees north (south negative) */
	Longitude float64     // Longitude, degrees east (west negative) */
	Lmst      float64     // Local mean sidereal time, degrees */
	Mnanom    float64     // Mean anomaly, degrees */
	Mnlong    float64     // Mean longitude, degrees */
	Rascen    float64     // Right ascension, degrees */
	Press     float64     // Surface pressure, millibars, used for refraction correction and ampress */
	Prime     float64     // Factor that normalizes Kt, Kn, etc. */
	Sbcf      float64     // Shadow-band correction factor */
	Sbwid     float64     // Shadow-band width (cm) */
	Sbrad     float64     // Shadow-band radius (cm) */
	Sbsky     float64     // Shadow-band sky factor */
	Solcon    float64     // Solar constant (NREL uses 1367 W/sq m) */
	Ssha      float64     // Sunset(/rise) hour angle, degrees */
	Sretr     float64     // Sunrise time, minutes from midnight, local, WITHOUT refraction */
	Ssetr     float64     // Sunset time, minutes from midnight, local, WITHOUT refraction */
	Temp      float64     // Ambient dry-bulb temperature, degrees C, used for refraction correction */
	Tilt      float64     // Degrees tilt from horizontal of panel */
	Timezone  float64     // Time zone, east (west negative). USA:  Mountain = -7, Central = -6, etc. */
	Tst       float64     // True solar time, minutes from midnight */
	Tstfix    float64     // True solar time - local standard time */
	Unprime   float64     // Factor that denormalizes Kt', Kn', etc. */
	Utime     float64     // Universal (Greenwich) standard time */
	Zenetr    float64     // Solar zenith angle, no atmospheric correction (= ETR) */
	Zenref    float64     // Solar zenith angle, deg. from zenith, refracted */
	Tdat      trigdata
}

func (sp *solpos) GetSunrise() time.Time {
	h, m, s := sp.calculateHourMinSec(sp.Sretr)
	dt := time.Date(sp.Year, time.Month(sp.Month), sp.Day, 0, 0, 0, 0, time.FixedZone("ManualTimeZone", int(sp.Timezone*3600)))
	return dt.Add(time.Hour*time.Duration(h) +
		time.Minute*time.Duration(m) +
		time.Second*time.Duration(s))

}
func (sp *solpos) calculateHourMinSec(decMinutes float64) (hours int, minutes int, seconds int) {
	hour := decMinutes / 60
	hours = int(math.Floor(hour))
	minutes = int(math.Floor(60 * (hour - float64(hours))))
	seconds = int((60.0 * (hour - float64(hours))) - float64(minutes)/60)
	if seconds < 0 {
		seconds = 0
	}
	return
}

func (sp *solpos) GetSunset() time.Time {
	h, m, s := sp.calculateHourMinSec(sp.Ssetr)
	dt := time.Date(sp.Year, time.Month(sp.Month), sp.Day, 0, 0, 0, 0, time.FixedZone("ManualTimeZone", int(sp.Timezone*3600)))
	return dt.Add(time.Hour*time.Duration(h) +
		time.Minute*time.Duration(m) +
		time.Second*time.Duration(s))
}

func (sp *solpos) Getdate() time.Time {
	return time.Date(sp.Year, time.Month(sp.Month), sp.Day, sp.Hour, sp.Minute, sp.Second, 0, time.FixedZone("ManualTimeZone", int(sp.Timezone*3600)))
}

func (sp *solpos) SetDate(dt time.Time) {
	_, offset := dt.Zone()
	sp.Year = dt.Year()
	sp.Month = int(dt.Month())
	sp.Day = dt.Day()
	sp.Daynum = dt.YearDay()
	sp.Hour = dt.Hour()
	sp.Minute = dt.Minute()
	sp.Second = dt.Second()
	sp.Timezone = float64(offset / 3600)
}

func (sp *solpos) SetDay(day int) {
	sp.Day = day
}

func (sp *solpos) SetDaynum(daynum int) {
	sp.Daynum = daynum
}

func (sp *solpos) SetFunction(function SPFunctions) {
	sp.Function = function
}

func (sp *solpos) SetHour(hour int) {
	sp.Hour = hour
}

func (sp *solpos) SetInterval(interval int) {
	sp.Interval = interval
}

func (sp *solpos) SetMinute(minute int) {
	sp.Minute = minute
}

func (sp *solpos) SetMonth(month int) {
	sp.Month = month
}

func (sp *solpos) SetSecond(second int) {
	sp.Second = second
}

func (sp *solpos) SetYear(year int) {
	sp.Year = year
}

func (sp *solpos) SetAspect(aspect float64) {
	sp.Aspect = aspect
}

func (sp *solpos) SetLatitude(latitude float64) {
	sp.Latitude = latitude
}

func (sp *solpos) SetLongitude(longitude float64) {
	sp.Longitude = longitude
}

func (sp *solpos) SetPress(press float64) {
	sp.Press = press
}

func (sp *solpos) SetSbwid(sbwid float64) {
	sp.Sbwid = sbwid
}

func (sp *solpos) SetSbrad(sbrad float64) {
	sp.Sbrad = sbrad
}

func (sp *solpos) SetSbsky(sbsky float64) {
	sp.Sbsky = sbsky
}

func (sp *solpos) SetSolcon(solcon float64) {
	sp.Solcon = solcon
}

func (sp *solpos) SetTemp(temp float64) {
	sp.Temp = temp
}

func (sp *solpos) SetTilt(tilt float64) {
	sp.Tilt = tilt
}

func (sp *solpos) SetTimezone(timezone float64) {
	sp.Timezone = timezone
}

func (sp *solpos) SetZenref(zenref float64) {
	sp.Zenref = zenref
}

func (sp *solpos) setTrigdata(tdat trigdata) {
	sp.Tdat = tdat
}
func (sp *solpos) GetDay() int {
	return sp.Day
}

func (sp *solpos) GetDaynum() int {
	return sp.Daynum
}

func (sp *solpos) GetFunction() SPFunctions {
	return sp.Function
}

func (sp *solpos) GetHour() int {
	return sp.Hour
}

func (sp *solpos) GetInterval() int {
	return sp.Interval
}

func (sp *solpos) GetMinute() int {
	return sp.Minute
}

func (sp *solpos) GetMonth() int {
	return sp.Month
}

func (sp *solpos) GetSecond() int {
	return sp.Second
}

func (sp *solpos) GetYear() int {
	return sp.Year
}

func (sp *solpos) GetAmass() float64 {
	return sp.Amass
}

func (sp *solpos) GetAmpress() float64 {
	return sp.Ampress
}

func (sp *solpos) GetAspect() float64 {
	return sp.Aspect
}

func (sp *solpos) GetAzim() float64 {
	return sp.Azim
}

func (sp *solpos) GetCosinc() float64 {
	return sp.Cosinc
}

func (sp *solpos) GetCoszen() float64 {
	return sp.Coszen
}

func (sp *solpos) GetDayang() float64 {
	return sp.Dayang
}

func (sp *solpos) GetDeclin() float64 {
	return sp.Declin
}

func (sp *solpos) GetEclong() float64 {
	return sp.Eclong
}

func (sp *solpos) GetEcobli() float64 {
	return sp.Ecobli
}

func (sp *solpos) GetEctime() float64 {
	return sp.Ectime
}

func (sp *solpos) GetElevetr() float64 {
	return sp.Elevetr
}

func (sp *solpos) GetElevref() float64 {
	return sp.Elevref
}

func (sp *solpos) GetEqntim() float64 {
	return sp.Eqntim
}

func (sp *solpos) GetErv() float64 {
	return sp.Erv
}

func (sp *solpos) GetEtr() float64 {
	return sp.Etr
}

func (sp *solpos) GetEtrn() float64 {
	return sp.Etrn
}

func (sp *solpos) GetEtrtilt() float64 {
	return sp.Etrtilt
}

func (sp *solpos) GetGmst() float64 {
	return sp.Gmst
}

func (sp *solpos) GetHrang() float64 {
	return sp.Hrang
}

func (sp *solpos) GetJulday() float64 {
	return sp.Julday
}

func (sp *solpos) GetLatitude() float64 {
	return sp.Latitude
}

func (sp *solpos) GetLongitude() float64 {
	return sp.Longitude
}

func (sp *solpos) GetLmst() float64 {
	return sp.Lmst
}

func (sp *solpos) GetMnanom() float64 {
	return sp.Mnanom
}

func (sp *solpos) GetMnlong() float64 {
	return sp.Mnlong
}

func (sp *solpos) GetRascen() float64 {
	return sp.Rascen
}

func (sp *solpos) GetPress() float64 {
	return sp.Press
}

func (sp *solpos) GetPrime() float64 {
	return sp.Prime
}

func (sp *solpos) GetSbcf() float64 {
	return sp.Sbcf
}

func (sp *solpos) GetSbwid() float64 {
	return sp.Sbwid
}

func (sp *solpos) GetSbrad() float64 {
	return sp.Sbrad
}

func (sp *solpos) GetSbsky() float64 {
	return sp.Sbsky
}

func (sp *solpos) GetSolcon() float64 {
	return sp.Solcon
}

func (sp *solpos) GetSsha() float64 {
	return sp.Ssha
}

func (sp *solpos) GetSretr() float64 {
	return sp.Sretr
}

func (sp *solpos) GetSsetr() float64 {
	return sp.Ssetr
}

func (sp *solpos) GetTemp() float64 {
	return sp.Temp
}

func (sp *solpos) GetTilt() float64 {
	return sp.Tilt
}

func (sp *solpos) GetTimezone() float64 {
	return sp.Timezone
}

func (sp *solpos) GetTst() float64 {
	return sp.Tst
}

func (sp *solpos) GetTstfix() float64 {
	return sp.Tstfix
}

func (sp *solpos) GetUnprime() float64 {
	return sp.Unprime
}

func (sp *solpos) GetUtime() float64 {
	return sp.Utime
}

func (sp *solpos) GetZenetr() float64 {
	return sp.Zenetr
}

func (sp *solpos) GetZenref() float64 {
	return sp.Zenref
}

/*============================================================================
*    Long int function S_solpos, adapted from the NREL VAX solar libraries
*
*    This function calculates the apparent solar position and intensity
*    (theoretical maximum solar energy) based on the date, time, and
*    location on Earth. (DEFAULT values are from the optional S_posinit
*    function.)
*
*    Requires:
*        Date and time:
*            year
*            month  (optional without daynum)
*            day    (optional without daynum)
*            daynum
*            hour
*            minute
*            second
*        Location:
*            latitude
*            longitude
*        Location/time adjuster:
*            timezone
*        Atmospheric pressure and temperature:
*            press     DEFAULT 1013.0 mb
*            temp      DEFAULT 10.0 degrees C
*        Tilt of flat surface that receives solar energy:
*            aspect    DEFAULT 180 (South)
*            tilt      DEFAULT 0 (Horizontal)
*        Shadow band parameters:
*            sbwid     DEFAULT 7.6 cm
*            sbrad     DEFAULT 31.7 cm
*            sbsky     DEFAULT 0.04
*        Functionality
*            function  DEFAULT S_ALL (all output parameters computed)
*
*    Returns:
*        everything defined at the top of this listing.
*----------------------------------------------------------------------------*/

func (sp *solpos) Calculate() error {
	// renew the date
	sp.SetDate(sp.Getdate())
	/* validate the inputs */
	err := sp.validate()
	if err != nil {
		return err
	}
	if sp.Function == 0 {
		return errors.New("No function set")
	}

	if sp.Function.HasFlag(LDoy) {
		/* convert input doy to month-day */
		sp.doy2dom()
	} else {
		/* convert input month-day to doy */
		sp.dom2doy()
	}

	if sp.Function.HasFlag(LGeom) {
		/* do basic geometry calculations */
		sp.geometry()
	}

	if sp.Function.HasFlag(LZenetr) {
		/* etr at non-refracted zenith angle */
		sp.zen_no_ref()
	}

	if sp.Function.HasFlag(LSsha) {
		/* Sunset hour calculation */
		sp.ssha()
	}

	if sp.Function.HasFlag(LSbcf) {
		/* Shadowband correction factor */
		sp.sbcf()
	}

	if sp.Function.HasFlag(LTst) {
		/* true solar time */
		sp.tst()
	}

	if sp.Function.HasFlag(LSrss) {
		/* sunrise/sunset calculations */
		sp.srss()
	}

	if sp.Function.HasFlag(LSolazm) {
		/* solar azimuth calculations */
		sp.sazm()
	}

	if sp.Function.HasFlag(LRefrac) {
		/* atmospheric refraction calculations */

		sp.refrac()
	}

	if sp.Function.HasFlag(LAmass) {

		/* airmass calculations */
		sp.amass()
	}

	if sp.Function.HasFlag(LPrime) {
		/* kt-prime/unprime calculations */
		sp.prime()
	}

	if sp.Function.HasFlag(LEtr) {
		/* ETR and ETRN (refracted) */
		sp.etr()
	}

	if sp.Function.HasFlag(LTilt) {
		/* tilt calculations */
		sp.tilt()
	}

	return nil
}

/*============================================================================
*    Void function S_init
*
*    This function initiates all of the input functions to S_Solpos().
*    NOTE: This function is optional if you initialize all input parameters
*          in your calling code.
*
*    Requires: Pointer to a posdata structure, members of which are
*           initialized.
*
*    Returns: Void
*
*----------------------------------------------------------------------------*/
func (sp *solpos) init() {
	sp.Day = -99              /* Day of month (May 27 = 27, etc.) */
	sp.Daynum = -999          /* Day number (day of year; Feb 1 = 32 ) */
	sp.Hour = -99             /* Hour of day, 0 - 23 */
	sp.Minute = -99           /* Minute of hour, 0 - 59 */
	sp.Month = -99            /* Month number (Jan = 1, Feb = 2, etc.) */
	sp.Second = -99           /* Second of minute, 0 - 59 */
	sp.Year = -99             /* 4-digit year */
	sp.Interval = 0           /* instantaneous measurement interval */
	sp.Aspect = 180.0         /* Azimuth of panel surface (direction it faces) N=0, E=90, S=180, W=270 */
	sp.Latitude = -99.0       /* Latitude, degrees north (south negative) */
	sp.Longitude = -999.0     /* Longitude, degrees east (west negative) */
	sp.Press = 1013.0         /* Surface pressure, millibars */
	sp.Solcon = 1367.0        /* Solar constant, 1367 W/sq m */
	sp.Temp = 15.0            /* Ambient dry-bulb temperature, degrees C */
	sp.Tilt = 0.0             /* Degrees tilt from horizontal of panel */
	sp.Timezone = -99.0       /* Time zone, east (west negative). */
	sp.Sbwid = 7.6            /* Eppley shadow band width */
	sp.Sbrad = 31.7           /* Eppley shadow band radius */
	sp.Sbsky = 0.04           /* Drummond factor for partly cloudy skies */
	sp.Function.AddFlag(SAll) /* compute all parameters */
}

/*++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 *
 * Structures defined for this module
 *
 *++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
type trigdata struct /* used to pass calculated values locally */
{
	Cd float64 /* cosine of the declination */
	Ch float64 /* cosine of the hour angle */
	Cl float64 /* cosine of the latitude */
	Sd float64 /* sine of the declination */
	Sl float64 /* sine of the latitude */
}

/*++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
 *
 * Temporary global variables used only in this file:
 *
 *++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++*/
var month_days = [2][13]int{{0, 0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334}, {0, 0, 31, 60, 91, 121, 152, 182, 213, 244, 274, 305, 335}}

/* cumulative number of days prior to beginning of month */
const degrad float64 = 57.295779513 /* converts from radians to degrees */
const raddeg float64 = 0.0174532925 /* converts from degrees to radians */

/*============================================================================
 *    Local function prototypes
 ============================================================================*/
func (sp *solpos) validate() error {

	/* No absurd dates, please. */
	if sp.Function.HasFlag(LGeom) {

		if (sp.Year < 1950) || (sp.Year > 2050) { /* limits of algoritm */

			return errors.New("Please fix the year: [1950-2050]")
		}
		if !(sp.Function.HasFlag(SDoy)) && ((sp.Month < 1) || (sp.Month > 12)) {
			return errors.New("Please fix the month [1-12]")
		}
		if !(sp.Function.HasFlag(SDoy)) && ((sp.Day < 1) || (sp.Day > 31)) {
			return errors.New("Please fix the day [1-31]")
		}
		if (sp.Function.HasFlag(SDoy)) && ((sp.Daynum < 1) || (sp.Daynum > 366)) {
			return errors.New("Please fix the day of year [1-366")
		}

		/* No absurd times, please. */
		if (sp.Hour < 0) || (sp.Hour > 24) {
			return errors.New("Please fix hour [0-24]")
		}
		if (sp.Minute < 0) || (sp.Minute > 59) {
			return errors.New("Please fix minute [0-59]")
		}
		if (sp.Second < 0) || (sp.Second > 59) {
			return errors.New("Please fix second [0-59]")
		}
		if (sp.Hour == 24) && (sp.Minute > 0) { /* no more than 24 hrs */

			return errors.New("Please fix hour and minute")
		}
		if (sp.Hour == 24) && (sp.Second > 0) { /* no more than 24 hrs */

			return errors.New("Please fix hour and second")
		}
		if math.Abs(sp.Timezone) > 12.0 {
			return errors.New("Please fix timezone [-12 - +12]")
		}
		if (sp.Interval < 0) || (sp.Interval > 28800) {
			return errors.New("Please fix interval (seconds) [0 - 28800]")
		}

		/* No absurd locations, please. */
		if math.Abs(sp.Longitude) > 180.0 {
			return errors.New("Please fix longitude [-180 - +180]")
		}
		if math.Abs(sp.Latitude) > 90.0 {
			return errors.New("Please fix latitude [-90 - +90]")
		}
		return nil
	}

	/* No silly temperatures or pressures, please. */
	if (sp.Function.HasFlag(LRefrac)) && (math.Abs(sp.Temp) > 100.0) {
		return errors.New("Please fix temperature [-100 - +100]")
	}

	if (sp.Function.HasFlag(LRefrac)) &&
		(sp.Press < 0.0) || (sp.Press > 2000.0) {
		return errors.New("Please fix press [0-2000]")
	}

	/* No out of bounds tilts, please */
	if (sp.Function.HasFlag(LTilt)) && (math.Abs(sp.Tilt) > 180.0) {
		return errors.New("Please fix tilt [-90 - 90]")
	}

	if (sp.Function.HasFlag(LTilt)) && (math.Abs(sp.Aspect) > 360.0) {
		return errors.New("Please fix aspect [-360 - 360]")
	}

	/* No oddball shadowbands, please */
	if (sp.Function.HasFlag(LSbcf)) &&
		(sp.Sbwid < 1.0) || (sp.Sbwid > 100.0) {
		return errors.New("Please fix shadow band width cm [1-100]")
	}

	if (sp.Function.HasFlag(LSbcf)) && (sp.Sbrad < 1.0) || (sp.Sbrad > 100.0) {
		return errors.New("Please fix shadow band radius (cm) [1-100]")
	}

	if (sp.Function.HasFlag(LSbcf)) && (math.Abs(sp.Sbsky) > 1.0) {
		return errors.New("Please fix shadow band sky factor [-1-+1]")
	}

	return nil
}

/*============================================================================
 *    Local Void function dom2doy
 *
 *    Converts day-of-month to day-of-year
 *
 *    Requires (from struct posdata parameter):
 *            year
 *            month
 *            day
 *
 *    Returns (via the struct posdata parameter):
 *            year
 *            daynum
 *----------------------------------------------------------------------------*/
func (sp *solpos) dom2doy() {
	sp.Daynum = sp.Day + month_days[0][sp.Month]

	/* (adjust for leap year) */
	if sp.Year%4 == 0 && (sp.Year%100 != 0 || sp.Year%400 == 0) && sp.Month > 2 {
		sp.Daynum += 1
	}

}

/*============================================================================
 *    Local void function doy2dom
 *
 *    This function computes the month/day from the day number.
 *
 *    Requires (from struct posdata parameter):
 *        Year and day number:
 *            year
 *            daynum
 *
 *    Returns (via the struct posdata parameter):
 *            year
 *            month
 *            day
 *----------------------------------------------------------------------------*/
func (sp *solpos) doy2dom() {
	var imon int /* Month (month_days) array counter */
	var leap int /* leap year switch */

	/* Set the leap year switch */
	if ((sp.Year % 4) == 0) &&
		(((sp.Year % 100) != 0) || ((sp.Year % 400) == 0)) {
		leap = 1
	} else {
		leap = 0
	}

	/* Find the month */
	imon = 12

	for sp.Daynum <= month_days[leap][imon] {
		imon--
	}

	/* Set the month and day of month */
	sp.Month = imon
	sp.Day = sp.Daynum - month_days[leap][imon]

}

/*============================================================================
 *    Local Void function geometry
 *
 *    Does the underlying geometry for a given time and location
 *----------------------------------------------------------------------------*/
func (sp *solpos) geometry() {
	var bottom float64 /* denominator (bottom) of the fraction */
	var c2 float64     /* cosine of d2 */
	var cd float64     /* cosine of the day angle or delination */
	var d2 float64     /* pdat->dayang times two */
	var delta float64  /* difference between current year and 1949 */
	var s2 float64     /* sine of d2 */
	var sd float64     /* sine of the day angle */
	var top float64    /* numerator (top) of the fraction */
	var leap int       /* leap year counter */

	/* Day angle */
	/*  Iqbal, M.  1983.  An Introduction to Solar Radiation.
	    Academic Press, NY., page 3 */
	sp.Dayang = 360.0 * (float64(sp.Daynum) - 1.0) / 365.0

	/* Earth radius vector * solar constant = solar energy */
	/*  Spencer, J. W.  1971.  Fourier series representation of the
	    position of the sun.  Search 2 (5), page 172 */
	sd = math.Sin(raddeg * sp.Dayang)
	cd = math.Cos(raddeg * sp.Dayang)
	d2 = 2.0 * sp.Dayang
	c2 = math.Cos(raddeg * d2)
	s2 = math.Sin(raddeg * d2)

	sp.Erv = 1.000110 + 0.034221*cd + 0.001280*sd
	sp.Erv += 0.000719*c2 + 0.000077*s2

	/* Universal Coordinated (Greenwich standard) time */
	/*  Michalsky, J.  1988.  The Astronomical Almanac's algorithm for
	    approximate solar position (1950-2050).  Solar Energy 40 (3),
	    pp. 227-235. */
	sp.Utime = float64(sp.Hour*3600.0 + sp.Minute*60.0 + sp.Second - sp.Interval/2.0)
	sp.Utime = sp.Utime/3600.0 - sp.Timezone

	/* Julian Day minus 2,400,000 days (to eliminate roundoff errors) */
	/*  Michalsky, J.  1988.  The Astronomical Almanac's algorithm for
	    approximate solar position (1950-2050).  Solar Energy 40 (3),
	    pp. 227-235. */

	/* No adjustment for century non-leap years since this function is
	   bounded by 1950 - 2050 */
	delta = float64(sp.Year - 1949)
	leap = (int((delta / 4.0)))
	sp.Julday = 32916.5 + (delta * 365.0) + float64(leap) + float64(sp.Daynum) + (sp.Utime / 24.0)

	/* Time used in the calculation of ecliptic coordinates */
	/* Noon 1 JAN 2000 = 2,400,000 + 51,545 days Julian Date */
	/*  Michalsky, J.  1988.  The Astronomical Almanac's algorithm for
	    approximate solar position (1950-2050).  Solar Energy 40 (3),
	    pp. 227-235. */
	sp.Ectime = sp.Julday - 51545.0

	/* Mean longitude */
	/*  Michalsky, J.  1988.  The Astronomical Almanac's algorithm for
	    approximate solar position (1950-2050).  Solar Energy 40 (3),
	    pp. 227-235. */
	sp.Mnlong = 280.460 + 0.9856474*sp.Ectime

	/* (dump the multiples of 360, so the answer is between 0 and 360) */
	sp.Mnlong -= float64(360.0 * int(sp.Mnlong/360.0))
	if sp.Mnlong < 0.0 {
		sp.Mnlong += 360.0
	}

	/* Mean anomaly */
	/*  Michalsky, J.  1988.  The Astronomical Almanac's algorithm for
	    approximate solar position (1950-2050).  Solar Energy 40 (3),
	    pp. 227-235. */
	sp.Mnanom = 357.528 + 0.9856003*sp.Ectime

	/* (dump the multiples of 360, so the answer is between 0 and 360) */
	sp.Mnanom -= float64(360.0 * int(sp.Mnanom/360.0))
	if sp.Mnanom < 0.0 {
		sp.Mnanom += 360.0
	}

	/* Ecliptic longitude */
	/*  Michalsky, J.  1988.  The Astronomical Almanac's algorithm for
	    approximate solar position (1950-2050).  Solar Energy 40 (3),
	    pp. 227-235. */
	sp.Eclong = sp.Mnlong + 1.915*math.Sin(sp.Mnanom*raddeg) + 0.020*math.Sin(2.0*sp.Mnanom*raddeg)

	/* (dump the multiples of 360, so the answer is between 0 and 360) */
	sp.Eclong -= float64(360 * int(sp.Eclong/360.0))
	if sp.Eclong < 0.0 {
		sp.Eclong += 360.0
	}

	/* Obliquity of the ecliptic */
	/*  Michalsky, J.  1988.  The Astronomical Almanac's algorithm for
	    approximate solar position (1950-2050).  Solar Energy 40 (3),
	    pp. 227-235. */

	/* 02 Feb 2001 SMW corrected sign in the following line */
	/*  pdat->ecobli = 23.439 + 4.0e-07 * pdat->ectime;     */
	sp.Ecobli = 23.439 - 4.0e-07*sp.Ectime

	/* Declination */
	/*  Michalsky, J.  1988.  The Astronomical Almanac's algorithm for
	    approximate solar position (1950-2050).  Solar Energy 40 (3),
	    pp. 227-235. */
	sp.Declin = degrad * math.Asin(math.Sin(sp.Ecobli*raddeg)*math.Sin(sp.Eclong*raddeg))

	/* Right ascension */
	/*  Michalsky, J.  1988.  The Astronomical Almanac's algorithm for
	    approximate solar position (1950-2050).  Solar Energy 40 (3),
	    pp. 227-235. */
	top = math.Cos(raddeg*sp.Ecobli) * math.Sin(raddeg*sp.Eclong)
	bottom = math.Cos(raddeg * sp.Eclong)

	sp.Rascen = degrad * math.Atan2(top, bottom)

	/* (make it a positive angle) */
	if sp.Rascen < 0.0 {
		sp.Rascen += 360.0
	}

	/* Greenwich mean sidereal time */
	/*  Michalsky, J.  1988.  The Astronomical Almanac's algorithm for
	    approximate solar position (1950-2050).  Solar Energy 40 (3),
	    pp. 227-235. */
	sp.Gmst = 6.697375 + 0.0657098242*sp.Ectime + sp.Utime

	/* (dump the multiples of 24, so the answer is between 0 and 24) */
	sp.Gmst -= float64(24 * (int(sp.Gmst / 24.0)))
	if sp.Gmst < 0.0 {
		sp.Gmst += 24.0
	}

	/* Local mean sidereal time */
	/*  Michalsky, J.  1988.  The Astronomical Almanac's algorithm for
	    approximate solar position (1950-2050).  Solar Energy 40 (3),
	    pp. 227-235. */
	sp.Lmst = sp.Gmst*15.0 + sp.Longitude

	/* (dump the multiples of 360, so the answer is between 0 and 360) */
	sp.Lmst -= float64(360 * (int(sp.Lmst / 360.0)))
	if sp.Lmst < 0.0 {
		sp.Lmst += 360.0
	}

	/* Hour angle */
	/*  Michalsky, J.  1988.  The Astronomical Almanac's algorithm for
	    approximate solar position (1950-2050).  Solar Energy 40 (3),
	    pp. 227-235. */
	sp.Hrang = sp.Lmst - sp.Rascen

	/* (force it between -180 and 180 degrees) */
	if sp.Hrang < -180.0 {
		sp.Hrang += 360.0
	}
	if sp.Hrang > 180.0 {
		sp.Hrang -= 360.0
	}

}

/*============================================================================
 *    Local Void function zen_no_ref
 *
 *    ETR solar zenith angle
 *       Iqbal, M.  1983.  An Introduction to Solar Radiation.
 *            Academic Press, NY., page 15
 *----------------------------------------------------------------------------*/
func (sp *solpos) zen_no_ref() {
	var cz float64 /* cosine of the solar zenith angle */

	sp.localtrig()
	cz = sp.Tdat.Sd*sp.Tdat.Sl + sp.Tdat.Cd*sp.Tdat.Cl*sp.Tdat.Ch

	/* (watch out for the roundoff errors) */
	if math.Abs(cz) > 1.0 {
		if cz >= 0.0 {
			cz = 1.0
		} else {
			cz = -1.0
		}

	}

	sp.Zenetr = math.Acos(cz) * degrad

	/* (limit the degrees below the horizon to 9 [+90 -> 99]) */
	if sp.Zenetr > 99.0 {
		sp.Zenetr = 99.0
	}

	sp.Elevetr = 90.0 - sp.Zenetr
}

/*============================================================================
 *    Local Void function ssha
 *
 *    Sunset hour angle, degrees
 *       Iqbal, M.  1983.  An Introduction to Solar Radiation.
 *            Academic Press, NY., page 16
 *----------------------------------------------------------------------------*/
func (sp *solpos) ssha() {
	var cssha float64 /* cosine of the sunset hour angle */
	var cdcl float64  /* ( cd * cl ) */

	sp.localtrig()
	cdcl = sp.Tdat.Cd * sp.Tdat.Cl

	if math.Abs(cdcl) >= 0.001 {
		cssha = -sp.Tdat.Sl * sp.Tdat.Sd / cdcl

		/* This keeps the cosine from blowing on roundoff */
		if cssha < -1.0 {
			sp.Ssha = 180.0
		} else if cssha > 1.0 {
			sp.Ssha = 0.0
		} else {
			sp.Ssha = degrad * math.Acos(cssha)
		}
	} else if ((sp.Declin >= 0.0) && (sp.Latitude > 0.0)) || ((sp.Declin < 0.0) && (sp.Latitude < 0.0)) {
		sp.Ssha = 180.0
	} else {
		sp.Ssha = 0.0
	}

}

/*============================================================================
 *    Local Void function sbcf
 *
 *    Shadowband correction factor
 *       Drummond, A. J.  1956.  A contribution to absolute pyrheliometry.
 *            Q. J. R. Meteorol. Soc. 82, pp. 481-493
 *----------------------------------------------------------------------------*/
func (sp *solpos) sbcf() {
	var p, t1, t2 float64 /* used to compute sbcf */
	sp.localtrig()
	p = 0.6366198 * sp.Sbwid / sp.Sbrad * math.Pow(sp.Tdat.Cd, 3)
	t1 = sp.Tdat.Sl * sp.Tdat.Sd * sp.Ssha * raddeg
	t2 = sp.Tdat.Cl * sp.Tdat.Cd * math.Sin(sp.Ssha*raddeg)
	sp.Sbcf = sp.Sbsky + 1.0/(1.0-p*(t1+t2))
}

/*============================================================================
 *    Local Void function tst
 *
 *    TST -> True Solar Time = local standard time + TSTfix, time
 *      in minutes from midnight.
 *        Iqbal, M.  1983.  An Introduction to Solar Radiation.
 *            Academic Press, NY., page 13
 *----------------------------------------------------------------------------*/
func (sp *solpos) tst() {
	sp.Tst = (180.0 + sp.Hrang) * 4.0
	sp.Tstfix = sp.Tst - float64(sp.Hour)*60.0 - float64(sp.Minute) - float64(sp.Second)/60.0 + float64(sp.Interval)/120.0 /* add back half of the interval */

	/* bound tstfix to this day */
	for sp.Tstfix > 720.0 {
		fmt.Println("bigger")
		sp.Tstfix -= 1440.0
	}

	for sp.Tstfix < -720.0 {
		fmt.Println("smaller")
		sp.Tstfix += 1440.0
	}

	sp.Eqntim = sp.Tstfix + 60.0*sp.Timezone - 4.0*sp.Longitude
}

/*============================================================================
 *    Local Void function srss
 *
 *    Sunrise and sunset times (minutes from midnight)
 *----------------------------------------------------------------------------*/
func (sp *solpos) srss() {
	if sp.Ssha <= 1.0 {
		sp.Sretr = 2999.0
		sp.Ssetr = -2999.0
	} else if sp.Ssha >= 179.0 {
		sp.Sretr = -2999.0
		sp.Ssetr = 2999.0
	} else {
		sp.Sretr = 720.0 - 4.0*sp.Ssha - sp.Tstfix
		sp.Ssetr = 720.0 + 4.0*sp.Ssha - sp.Tstfix
	}
}

/*============================================================================
 *    Local Void function sazm
 *
 *    Solar azimuth angle
 *       Iqbal, M.  1983.  An Introduction to Solar Radiation.
 *            Academic Press, NY., page 15
 *----------------------------------------------------------------------------*/
func (sp *solpos) sazm() {
	var ca float64   /* cosine of the solar azimuth angle */
	var ce float64   /* cosine of the solar elevation */
	var cecl float64 /* ( ce * cl ) */
	var se float64   /* sine of the solar elevation */

	sp.localtrig()
	ce = math.Cos(raddeg * sp.Elevetr)
	se = math.Sin(raddeg * sp.Elevetr)

	sp.Azim = 180.0
	cecl = ce * sp.Tdat.Cl
	if math.Abs(cecl) >= 0.001 {
		ca = (se*sp.Tdat.Sl - sp.Tdat.Sd) / cecl
		if ca > 1.0 {
			ca = 1.0
		} else if ca < -1.0 {
			ca = -1.0
		}

		sp.Azim = 180.0 - math.Acos(ca)*degrad
		if sp.Hrang > 0 {
			sp.Azim = 360.0 - sp.Azim
		}
	}
}

/*============================================================================
 *    Local Int function refrac
 *
 *    Refraction correction, degrees
 *        Zimmerman, John C.  1981.  Sun-pointing programs and their
 *            accuracy.
 *            SAND81-0761, Experimental Systems Operation Division 4721,
 *            Sandia National Laboratories, Albuquerque, NM.
 *----------------------------------------------------------------------------*/
func (sp *solpos) refrac() {
	var prestemp float64 /* temporary pressure/temperature correction */
	var refcor float64   /* temporary refraction correction */
	var tanelev float64  /* tangent of the solar elevation angle */

	/* If the sun is near zenith, the algorithm bombs; refraction near 0 */
	if sp.Elevetr > 85.0 {
		refcor = 0.0
	} else {
		/* Otherwise, we have refraction */
		tanelev = math.Tan(raddeg * sp.Elevetr)
		if sp.Elevetr >= 5.0 {
			refcor = 58.1/tanelev - 0.07/(math.Pow(tanelev, 3)) + 0.000086/(math.Pow(tanelev, 5))
		} else if sp.Elevetr >= -0.575 {
			refcor = 1735.0 + sp.Elevetr*(-518.2+sp.Elevetr*(103.4+sp.Elevetr*(-12.79+sp.Elevetr*0.711)))
		} else {
			refcor = -20.774 / tanelev
		}
		prestemp =
			(sp.Press * 283.0) / (1013.0 * (273.0 + sp.Temp))
		refcor *= prestemp / 3600.0

	}

	/* Refracted solar elevation angle */
	sp.Elevref = sp.Elevetr + refcor

	/* (limit the degrees below the horizon to 9) */
	if sp.Elevref < -9.0 {
		sp.Elevref = -9.0
	}

	/* Refracted solar zenith angle */
	sp.Zenref = 90.0 - sp.Elevref
	sp.Coszen = math.Cos(raddeg * sp.Zenref)
}
func (sp *solpos) amass() {
	if sp.Zenref > 93.0 {
		sp.Amass = -1.0
		sp.Ampress = -1.0
	} else {
		sp.Amass =
			1.0 / (math.Cos(raddeg*sp.Zenref) + (0.50572 *
				math.Pow(96.07995-sp.Zenref, -1.6364)))

		sp.Ampress = sp.Amass * sp.Press / 1013.0
	}
}

/*============================================================================
 *    Local Void function prime
 *
 *    Prime and Unprime
 *    Prime  converts Kt to normalized Kt', etc.
 *       Unprime deconverts Kt' to Kt, etc.
 *            Perez, R., P. Ineichen, Seals, R., & Zelenka, A.  1990.  Making
 *            full use of the clearness index for parameterizing hourly
 *            insolation conditions. Solar Energy 45 (2), pp. 111-114
 *----------------------------------------------------------------------------*/
func (sp *solpos) prime() {
	sp.Unprime = 1.031*math.Exp(-1.4/(0.9+9.4/sp.Amass)) + 0.1
	sp.Prime = 1.0 / sp.Unprime
}

/*============================================================================
 *    Local Void function etr
 *
 *    Extraterrestrial (top-of-atmosphere) solar irradiance
 *----------------------------------------------------------------------------*/
func (sp *solpos) etr() {
	if sp.Coszen > 0.0 {
		sp.Etrn = sp.Solcon * sp.Erv
		sp.Etr = sp.Etrn * sp.Coszen

	} else {
		sp.Etrn = 0.0
		sp.Etr = 0.0
	}
}

/*============================================================================
 *    Local Void function tilt
 *
 *    ETR on a tilted surface
 *----------------------------------------------------------------------------*/
func (sp *solpos) tilt() {
	var ca float64  /* cosine of the solar azimuth angle */
	var cp float64  /* cosine of the panel aspect */
	var ct float64  /* cosine of the panel tilt */
	var sa float64  /* sine of the solar azimuth angle */
	var spp float64 /* sine of the panel aspect */
	var st float64  /* sine of the panel tilt */
	var sz float64  /* sine of the refraction corrected solar zenith angle */

	/* Cosine of the angle between the sun and a tipped flat surface,
	   useful for calculating solar energy on tilted surfaces */
	ca = math.Cos(raddeg * sp.Azim)
	cp = math.Cos(raddeg * sp.Aspect)
	ct = math.Cos(raddeg * sp.Tilt)
	sa = math.Sin(raddeg * sp.Azim)
	spp = math.Sin(raddeg * sp.Aspect)
	st = math.Sin(raddeg * sp.Tilt)
	sz = math.Sin(raddeg * sp.Zenref)
	sp.Cosinc = sp.Coszen*ct + sz*st*(ca*cp+sa*spp)

	if sp.Cosinc > 0.0 {
		sp.Etrtilt = sp.Etrn * sp.Cosinc
	} else {
		sp.Etrtilt = 0.0
	}

}

/*============================================================================
 *    Local Void function localtrig
 *
 *    Does trig on internal variable used by several functions
 *----------------------------------------------------------------------------*/
func (sp *solpos) localtrig() {
	/* define masks to prevent calculation of uninitialized variables */

	if sp.Tdat.Sd < -900.0 { // sd was initialized -999 as flag
		sp.Tdat.Sd = 1.0 // reflag as having completed calculations
		if sp.Function.HasFlag(CdMask) {
			sp.Tdat.Cd = math.Cos(raddeg * sp.Declin)

			if sp.Function.HasFlag(ChMask) {
			}
			sp.Tdat.Ch = math.Cos(raddeg * sp.Hrang)

			if sp.Function.HasFlag(ClMask) {
			}
			sp.Tdat.Cl = math.Cos(raddeg * sp.Latitude)

			if sp.Function.HasFlag(SdMask) {
			}
			sp.Tdat.Sd = math.Sin(raddeg * sp.Declin)

			if sp.Function.HasFlag(SlMask) {
			}
			sp.Tdat.Sl = math.Sin(raddeg * sp.Latitude)
		}

	}
}
