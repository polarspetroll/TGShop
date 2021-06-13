module telegramshop

go 1.15

require (
  app v0.0.0
  telegram v0.0.0
)


replace (
  app => ./app
  telegram => ./telegram
  DB => ./DB
)
