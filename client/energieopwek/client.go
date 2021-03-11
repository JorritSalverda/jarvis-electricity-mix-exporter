package energieopwek

type Client interface {
}

func NewClient() (Client, error) {
	return &client{}, nil
}

type client struct {
}

// def get_production_data_energieopwek(date, session=None):
//     r = session or requests.session()

//     # The API returns values per day from local time midnight until the last
//     # round 10 minutes if the requested date is today or for the entire day if
//     # it's in the past. 'sid' can be anything.
//     url = 'http://energieopwek.nl/jsonData.php?sid=2ecde3&Day=%s' % date.format('YYYY-MM-DD')
//     response = r.get(url)
//     obj = response.json()
//     production_input = obj['TenMin']['Country']

//     # extract the power values in kW from the different production types
//     # we only need column 0, 1 and 3 contain energy sum values
//     df_solar =    pd.DataFrame(production_input['Solar'])       .drop(['1','3'], axis=1).astype(int).rename(columns={"0" : "solar"})
//     df_offshore = pd.DataFrame(production_input['WindOffshore']).drop(['1','3'], axis=1).astype(int)
//     df_onshore =  pd.DataFrame(production_input['Wind'])        .drop(['1','3'], axis=1).astype(int)

//     # We don't differentiate between onshore and offshore wind so we sum them
//     # toghether and build a single data frame with named columns
//     df_wind = df_onshore.add(df_offshore).rename(columns={"0": "wind"})
//     df = pd.concat([df_solar, df_wind], axis=1)

//     # resample from 10min resolution to 15min resolution to align with ENTSOE data
//     # we duplicate every row and then group them per 3 and take the mean
//     df = pd.concat([df]*2).sort_index(axis=0).reset_index(drop=True).groupby(by=lambda x : math.floor(x/3)).mean()

//     # Convert kW to MW with kW resolution
//     df = df.apply(lambda x: round(x / 1000, 3))

//     return df
