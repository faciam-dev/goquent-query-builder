CREATE TABLE profiles (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT,
    age INT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
