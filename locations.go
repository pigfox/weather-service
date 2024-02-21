package main

import "weather-service/structs"

func getLocations() structs.Locations {
	locations := structs.Locations{}

	location := structs.Location{
		Lat:  "34.058032451466175",
		Lon:  "-118.23595125080323",
		Name: "Los Angeles",
	}
	locations.Area = append(locations.Area, location)

	location = structs.Location{
		Lat:  "37.77368416198036",
		Lon:  "-122.41185482459379",
		Name: "San Francisco",
	}
	locations.Area = append(locations.Area, location)

	location = structs.Location{
		Lat:  "32.71622791446506",
		Lon:  "-117.16627851159686",
		Name: "San Diego",
	}
	locations.Area = append(locations.Area, location)

	location = structs.Location{
		Lat:  "61.217214806510825",
		Lon:  "-149.86891662783154",
		Name: "Anchorage",
	}
	locations.Area = append(locations.Area, location)

	location = structs.Location{
		Lat:  "21.309463585349853",
		Lon:  "-157.8570993659206",
		Name: "Honolulu",
	}
	locations.Area = append(locations.Area, location)

	location = structs.Location{
		Lat:  "-33.86955210512373",
		Lon:  "151.2079001460676",
		Name: "Sydney",
	}
	locations.Area = append(locations.Area, location)

	location = structs.Location{
		Lat:  "-34.60449337243565",
		Lon:  "-58.341978798543366",
		Name: "Buenos Aires",
	}
	locations.Area = append(locations.Area, location)

	location = structs.Location{
		Lat:  "-33.92546155535883",
		Lon:  "18.432284113203856",
		Name: "Cape Town",
	}
	locations.Area = append(locations.Area, location)

	location = structs.Location{
		Lat:  "55.676326942940015",
		Lon:  "12.563113264981702",
		Name: "Copenhagen",
	}
	locations.Area = append(locations.Area, location)

	location = structs.Location{
		Lat:  "55.95367676508114",
		Lon:  "-3.1916086678208706",
		Name: "Edinburgh",
	}
	locations.Area = append(locations.Area, location)

	location = structs.Location{
		Lat:  "23.971261177497123",
		Lon:  "90.37722903083214",
		Name: "Dhaka",
	}
	locations.Area = append(locations.Area, location)

	location = structs.Location{
		Lat:  "-9.38938063045734",
		Lon:  "147.29739503837453",
		Name: "Port Moresby",
	}
	locations.Area = append(locations.Area, location)

	location = structs.Location{
		Lat:  "14.741560348469392",
		Lon:  "121.01875424393413",
		Name: "Manila",
	}
	locations.Area = append(locations.Area, location)

	return locations
}
