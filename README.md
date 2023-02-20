Project 1 - Magic internet money 🧙

Overview:
The project is a simple game platform called “AYO” where users can pay to play a game with 
either on-chain bitcoin payment or off-chain lightning. The platform will be built using 
Golang for the back-end and Vue.js with TypeScript for the front-end.

System Architecture:
The system will be comprised of the following components:

User Flow:
User selects “Lightning Payment” or “Bitcoin Payment” to make payment.
User is prompted to pay the required fee to start the game.
Once the payment is confirmed, the user’s wallet is funded.
User can initiate a game, and get debited from his wallet for starting a game.
The game is played between the two players, or against the computer.
On completion of the game, the winner of the game is declared.
The winner is rewarded with the payment made by the loser (user-to-user), OR refunded 50% of 
the game fee (user-to-computer).

The front-end, which will be built using Vue.js with TypeScript and will handle user 
interaction with the platform, Tailwind css for styling.
Views/Pages include:
SSO with Google
Dashboard - with button to start a game
Game board
Make Payment (Fund Wallet) - either with bitcoin/lightning
Payment history
Game history

The back-end, which will be built using Golang and will handle the game logic and payment 
processing. API endpoints include:
[GET] SSO auth callback
[GET] Payment history
[GET] Show payment details
[GET] Game history
[GET] Show game details
[POST] Start a game (auto debit)
[POST] Resign from a game
[POST] Initialize a payment
Generate a new address (for on-chain payment)
Create a Payment Invoice (for off-chaing lightning payment)
[POST] Sign out

A bitcoin wallet, which will be used to handle the bitcoin payments.

A lightning network node, which will be used to handle off-chain lightning payments.

Game Logic:
The game will be a simple version of the traditional African game called “AYO”. The game is a 
two-player game where each player takes turns moving the pieces. The game will be played on a 
board with 2x6 pots, with each player starting with 24 pieces on their side of the board. The 
objective of the game is to capture all of your opponent’s pieces.

Payment Processing:
The platform will allow users to fund their wallet with either on-chain bitcoin payment or 
lightning payment, after which they can pay from their wallet to play the game. The platform 
will handle payment processing using either the bitcoin wallet or the lightning network node. 
When a user initiates a game, they will be prompted to select their preferred payment method.

MVP Milestone:
API Completed
Play against a machine

Non-Goals:
There are no intentions of sweeping the funds out of the wallet(s).
There are no intentions to allow for fund withdrawal by the user.
There are no intentions to give customers full custody to the wallet by means of private key 
or seed phrase.



Introspection:
The concept adopted is to create a “virtual wallet” for a customer so as to reduce the hassle 
of payment for recurring game users. By so doing, a user can pay once, while his balance is 
being managed internally.
SSO with Google to be used to reduce the hassle of authentication and user identification.

Dependencies:
Vue 3 / Svelte
Tailwind css
Go
Btcd - github.com/btcsuite/btcd
Lndclient - github.com/lightninglabs/lndclient

Timeline:
Day 1:
Project setup
Complete SSO with Google OAuth2

Day 2:
Game design
Game play against computer.

Day 3:
Generate Bitcoin Address for deposit
Update wallet balance on payment confirmation

Day 4:
Create a Lightning Invoice; Update wallet balance on payment confirmation

Day 5:
Gate the gameplay with wallet debit, Wallet refunds

Day 6: 
Cleanup

