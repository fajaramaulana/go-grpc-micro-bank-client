package exception

import "github.com/rs/zerolog/log"

func PanicIfNeeded(err interface{}) {
	if err != nil {
		log.Panic().Msgf("%v", err)
		panic(err)
	}
}
