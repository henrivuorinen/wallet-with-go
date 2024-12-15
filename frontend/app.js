// Slot Machine Frontend - JavaScript Implementation

const appState = {
    balance: 0,
    username: '',
    isLoggedIn: false,
};

const DOM = {
    loginSection: document.getElementById('login-section'),
    registerSection: document.getElementById('register-section'),
    gameSection: document.getElementById('game-section'),
    balance: document.getElementById('balance'),
    usernameDisplay: document.getElementById('user-name'),
    message: document.getElementById('message'),
    wheels: [
        document.getElementById('wheel1'),
        document.getElementById('wheel2'),
        document.getElementById('wheel3'),
    ],
};

const registerPlayer = async () => {
    const response = await fetch('http://localhost:8080/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ player_id: 'player1', name: 'John Doe', balance: 100 })
    });
    const data = await response.json();
    console.log(data);
};

const loginPlayer = async () => {
    const response = await fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ player_id: 'player1' })
    });
    const data = await response.json();
    console.log(data);
};


// Utility Functions
function setView(view) {
    const sections = [DOM.loginSection, DOM.registerSection, DOM.gameSection];
    sections.forEach((section) => (section.style.display = 'none'));
    view.style.display = 'block';
}

function updateBalanceDisplay() {
    DOM.balance.textContent = appState.balance;
}

function showMessage(message) {
    DOM.message.textContent = message;
}

// Event Handlers
async function handleLogin(event) {
    event.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    if (!username || !password) {
        showMessage('Username and password are required.', true);
        return;
    }

    try {
        const response = await fetch('http://localhost:8080/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                player_id: username,
                password: password, // Include password in the API call
            }),
            credentials: "include",
            // Add this option to bypass SSL errors (development only):
            mode: "no-cors",
        });

        const data = await response.json();

        if (response.ok) {
            appState.isLoggedIn = true;
            appState.username = username;
            appState.balance = data.balance;

            DOM.usernameDisplay.textContent = username;
            updateBalanceDisplay();
            showMessage('Login successful!');
            setView(DOM.gameSection);
        } else {
            showMessage(`Login failed: ${data.error || 'Invalid credentials'}`, true);
        }
    } catch (error) {
        console.error('Error during login:', error);
        showMessage('Network error. Please try again later.', true);
    }
}


async function handleRegister(event) {
    event.preventDefault();
    const username = document.getElementById('new-username').value;
    const password = document.getElementById('new-password').value;

    if (!username || !password) {
        showMessage('Username and password are required.', true);
        return;
    }

    try {
        const response = await fetch('http://localhost:8080/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                player_id: username,
                password: password, // Include password in the API call
                name: username, // For simplicity, using username as the name
                balance: 100, // Initial balance
            }),
            credentials: "include",
            // Add this option to bypass SSL errors (development only):
            mode: "no-cors",
        });

        const data = await response.json();

        if (response.ok) {
            showMessage('Registration successful! Please log in.');
            setView(DOM.loginSection);
        } else {
            showMessage(`Registration failed: ${data.error || 'Unknown error'}`, true);
        }
    } catch (error) {
        console.error('Error during registration:', error);
        showMessage('Network error. Please try again later.', true);
    }
}


function handlePurchase() {
    appState.balance += 10; // Add balance locally
    updateBalanceDisplay();
}

async function handleSpin() {
    if (appState.balance < 1) {
        showMessage('Not enough balance to play!');
        return;
    }

    appState.balance--;
    updateBalanceDisplay();

    const symbols = ['🔵', '🔺', '⭐', '⬛'];
    const results = DOM.wheels.map(
        (wheel) => symbols[Math.floor(Math.random() * symbols.length)]
    );

    // Display results in the wheels
    results.forEach((symbol, index) => {
        DOM.wheels[index].textContent = symbol;
    });

    console.log('Spin results:', results);

    // Check if all symbols match
    const isWin = results.every((symbol) => symbol === results[0]);
    console.log('Win condition:', isWin);

    if (isWin) {
        const winnings = 50;
        if (!appState.transaction_id) appState.transaction_id = 1;
        const transaction_id = appState.transaction_id++;
        console.log('Calling /win API with:', {
            player_id: appState.username,
            transaction_id,
            amount: winnings,
        });

        try {
            const winResponse = await fetch('http://localhost:8080/win', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    //'X-Game-Engine-Api-Key': 'default-secure-api-key',
                },
                body: JSON.stringify({
                    player_id: appState.username,
                    transaction_id,
                    amount: winnings,
                }),
                credentials: "include",
                // Add this option to bypass SSL errors (development only):
                mode: "no-cors",
            });

            if (winResponse.ok) {
                appState.balance += winnings;
                updateBalanceDisplay();
                showMessage(`You win! Added ${winnings} to your balance.`);
                console.log('Win API call successful.');
            } else {
                showMessage(
                    'Error updating winnings. Please contact maintenance, they probably won’t answer though.'
                );
                console.error('Win API error:', winResponse.status, winResponse.statusText);
            }
        } catch (error) {
            console.error('Error during API call:', error);
            showMessage('Network error. Please try again later.');
        }
    } else {
        showMessage('Try again!');
    }
}

// Add Event Listeners
document.getElementById('login-form').addEventListener('submit', handleLogin);
document.getElementById('register-form').addEventListener('submit', handleRegister);
document.getElementById('purchase-balance').addEventListener('click', handlePurchase);
document.getElementById('spin').addEventListener('click', handleSpin);
document.getElementById('show-register').addEventListener('click', () => setView(DOM.registerSection));
document.getElementById('show-login').addEventListener('click', () => setView(DOM.loginSection));

// Initial Setup
setView(DOM.loginSection);
