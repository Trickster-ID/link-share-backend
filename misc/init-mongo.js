// Create or switch to a new database
db = db.getSiblingDB('link_share');

// Create a new collection
db.createCollection('refresh_token_sessions');

// Create indexes for refresh_token and exp auto delete field
db.refresh_token_sessions.createIndex({ refresh_token: 1 }, { unique: true });
db.refresh_token_sessions.createIndex({ exp: 1 },{ expireAfterSeconds: 0 });


// Create a new collection
db.createCollection('access_token_sessions');

// Create indexes for refresh_token and exp auto delete field
db.access_token_sessions.createIndex({ access_token: 1 }, { unique: true });
db.access_token_sessions.createIndex({ exp: 1 },{ expireAfterSeconds: 0 });