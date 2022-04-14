module github.com/MKwann7/zgEXCELL-3-Media/app/media

go 1.18

require (
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.4.0
	github.com/urfave/negroni v1.0.0
	gopkg.in/gographics/imagick.v2 v2.6.0
)

replace github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/controllers => ../media/src/code/controllers

replace github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/dtos => ../media/src/code/dtos

replace github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/libraries/builder => ../media/src/code/libraries/builder

replace github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/libraries/db => ../media/src/code/libraries/db

replace github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/libraries/exceptions => ../media/src/code/libraries/exceptions

replace github.com/MKwann7/zgEXCELL-3-Media/app/media/src/code/libraries/helper => ../media/src/code/libraries/helper
