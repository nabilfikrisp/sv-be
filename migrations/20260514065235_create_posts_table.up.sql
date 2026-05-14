CREATE TABLE IF NOT EXISTS posts (
    id INT NOT NULL AUTO_INCREMENT,
    title VARCHAR(200),
    content TEXT,
    category VARCHAR(100),
    created_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    status ENUM('publish', 'draft', 'thrash') DEFAULT 'draft',
    PRIMARY KEY (id),
    INDEX idx_status (status),
    INDEX idx_category (category),
    INDEX idx_created_date (created_date)
)