## now im in the services layer 

	## it would be implement game_service.go
	## which will use the player and game store
	## first functionality that the game store should have is the update board (execute move)
	## other things i dont know , but the board is the most important on






## AI TODO LIST< NOT 100 CORRECT>

Tic-Tac-Toe Backend: Week 2 Requirements
Phase 1: The Service Layer (The Brain)
[ ] Create internal/service/game_service.go and GameService struct.

[ ] Implement hydrateBoard and dehydrateBoard (Helper to turn DB string --- into [9]string).

[ ] Write CheckWinner logic (Scan the 8 winning patterns).

[ ] Write IsDraw logic (Check if board is full with no winner).

[ ] Implement ExecuteMove:

[ ] Validate it's the player's turn.

[ ] Validate the square is empty.

[ ] Update the board array.

[ ] Check for Win/Draw.

[ ] Call Store.UpdateMove (to save and flip the turn).

Phase 2: Authentication & Rooms
[ ] Create Login/Signup Handlers (Store users in Players table).

[ ] Implement Session Management (So the backend knows who is making the move).

[ ] Implement CreateGame (Generate unique room ID, set player X).

[ ] Implement JoinGame (Add player O, ensure only 2 unique users per game).

[ ] Create GetAllRooms endpoint (To show available games in the React UI).

Phase 3: Integration (The Handshake)
[ ] Set up CORS Middleware (Allow React localhost:3000 to talk to Go localhost:8080).

[ ] Create HTTP Handlers for Move, Join, and Create.

[ ] Implement Live Status logic (Return whose turn it is and the score to the frontend).

[ ] Final Test: Play a full game from start to finish using the React UI.
