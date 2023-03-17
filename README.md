# Project 1 - Magic internet money 🧙

### Overview:
The project is a simple game platform called “AYO” where users can pay to play a game with 
either on-chain bitcoin payment or off-chain lightning. The platform will be built using 
Golang for the back-end and Vue.js with TypeScript for the front-end.

The Design document of this project can be found [here](https://docs.google.com/document/d/1x30vNSXeKwguNqmqPlquS1Yltq1l3h9jhm4lhAinJqA/edit?usp=sharing)

Postman collection for the backend services can be found [here](https://documenter.getpostman.com/view/1439952/2s93Jxt2M6)

### Methodology
#### Payment Channels
This platform works with 2 channels,
1. Channel 1 - for making payments
2. Channel 2 - for receiving payments

Consequently, only the channel for making payments (channel 1) is meant to be locally funded. Channel 2
(which is used for receiving payments) doesn't have to be funded.

Gamers can pay generate and pay invoices to have access to a game.

#### Game Play
A gamer can only have one active game at each given time. Any other payments made while the gamer has an
active game is deposited into the internal wallet of the gamer.

#### Internal wallet
The internal wallet is a simple mechanism for improved user experience. 
Because the gamer has to pay for every game, it would be easier if the gamer gets to pay from his internal 
wallet, and (perhaps) fund the wallet once for extended use (over lightning).

### Get Started
To get started
1. Clone the repo <br/>
   `git clone git@github.com:Adesubomi/magic-ayo-api.git`
<br/><br/>
2. cd into project directory <br/>
   `cd magic-ayo-api`
<br/><br/>
3. Fetch go dependencies <br/>
    `go mod tidy`
<br/><br/>
4. Copy and update configuration file <br/>
    `cp config.example.toml config.toml`
<br/><br/>
5. Start the server <br/>
    `go run ./cmd/main.go --config=config.toml`
<br/><br/>

🎉 Let the games start!

