package main

func (s *db_store) addSPDropTables() {
	query := `DROP PROCEDURE IF EXISTS drop_all_tables ;`
	e := s.execSQLText(query)
	if e != nil {
		s.Panic(e)
	}
	query = `
	CREATE PROCEDURE drop_all_tables()
	BEGIN
	DECLARE done INT DEFAULT FALSE;
	DECLARE tableName VARCHAR(255);
	DECLARE cursor_tables CURSOR FOR
		SELECT table_name FROM information_schema.tables 
		WHERE table_type = 'BASE TABLE' AND table_schema = DATABASE();
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
	OPEN cursor_tables;
	read_loop: LOOP
		FETCH cursor_tables INTO tableName;
		IF done THEN
		LEAVE read_loop;
		END IF;
		SET @drop_sql = CONCAT('DROP TABLE IF EXISTS ', tableName, ';');
		PREPARE stmt_drop FROM @drop_sql;
		EXECUTE stmt_drop;
		DEALLOCATE PREPARE stmt_drop;
	END LOOP;
	CLOSE cursor_tables;
	END;
	`
	e = s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}
}

func (s *db_store) addSPDropTable() {
	query := `DROP PROCEDURE IF EXISTS drop_table ;`
	e := s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}
	query = `
	CREATE PROCEDURE drop_table(IN table_name VARCHAR(255))
	BEGIN
		SET @drop_sql = CONCAT('DROP TABLE IF EXISTS ', table_name, ';');

		PREPARE stmt_drop FROM @drop_sql;
		EXECUTE stmt_drop;
		DEALLOCATE PREPARE stmt_drop;
	END;
	`
	e = s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}
}

func (s *db_store) addSPSearchTables() {
	query := `DROP PROCEDURE IF EXISTS search_tables ;`
	e := s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}
	query = `
	CREATE PROCEDURE search_tables(IN pattern VARCHAR(255))
	BEGIN
		SET @sql = CONCAT('SELECT table_name FROM information_schema.tables WHERE table_name LIKE ''%', pattern, '%''');
		PREPARE stmt FROM @sql;
		EXECUTE stmt;
		DEALLOCATE PREPARE stmt;
	END ;
	`
	e = s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}
}

// 關聯
// 1、代理=>玩家清單
func (s *db_store) addSPUser() {
	query := `DROP PROCEDURE IF EXISTS generate_user;`
	e := s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}
	query = `
    CREATE PROCEDURE generate_user (IN agent_id VARCHAR(36))
    BEGIN
		SET @table_name = CONCAT('agent_' , agent_id) ;
        SET @sql = CONCAT(
            'CREATE TABLE IF NOT EXISTS ', @table_name, '_user (',
				'ID INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,',
				'LastIP VARCHAR(15) DEFAULT \'\',',
				'TheirUID VARCHAR(36) NOT NULL,',
				'TheirUName VARCHAR(20) NOT NULL,',
				'TheirUGrant VARCHAR(2) DEFAULT \'1\',',
				'Wallet DECIMAL(20, 4) DEFAULT \'0\',',
				'CreatedAt DATETIME NOT NULL,',
				'ModifiedAt DATETIME NOT NULL'
            ')'
        );
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END ;
	`
	e = s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}
}

// 關聯
// 1、代理=>玩家錢包
func (s *db_store) addSPWallet() {
	query := `DROP PROCEDURE IF EXISTS generate_wallet;`
	e := s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}
	query = `
    CREATE PROCEDURE generate_wallet (IN agent_id VARCHAR(36))
    BEGIN
		SET @table_name = CONCAT('agent_' , agent_id) ;
        SET @sql = CONCAT(
				'CREATE TABLE IF NOT EXISTS ', @table_name, '_wallet (',
				'ID INT UNSIGNED,',                            /** 我方 UID */
				'TheirUID VARCHAR(36) NOT NULL,',              /** 對方 UID */
				'GameID VARCHAR(4),',
				'RoundID VARCHAR(16),',
				'BeforeDiff DECIMAL(20, 4),',                  /** 異動前 */
				'Diff DECIMAL(20, 4),',                        /** 異動值 */
				'AfterDiff DECIMAL(20, 4),',                   /** 異動後 */
				'ActionType TINYINT UNSIGNED,',                /** <加/減> {0:支出,1:支出回滾,2:收入} */
				'TransactionDatetime DATETIME'                 /** UTC 異動時間 */
            ')'
        );
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END ;
	`
	e = s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}
}

func (s *db_store) addSPGame() {
	query := `DROP PROCEDURE IF EXISTS generate_game;`
	e := s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}
	query = `
    CREATE PROCEDURE generate_game (IN game_id VARCHAR(4))
    BEGIN
		SET @table_name = CONCAT('game_' , game_id, '_report') ;
        SET @sql = CONCAT(
            'CREATE TABLE IF NOT EXISTS ', @table_name, ' (',
				'RoundID VARCHAR(15) COMMENT "局號",',
				'TotalBet DECIMAL(20, 4),',
				'TotalWin DECIMAL(20, 4),',
				'Result TEXT COMMENT "回合結果",',
				'StartAt DATETIME NOT NULL,',
				'EndAt DATETIME NOT NULL'
            ')'
        );
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;

		SET @table_name = CONCAT('game_' , game_id, '_detail') ;
		SET @sql = CONCAT(
            'CREATE TABLE IF NOT EXISTS ', @table_name, ' (',
				'ID INT UNSIGNED NOT NULL COMMENT "玩家 ID",',
				'Serial VARCHAR(15) NOT NULL COMMENT "局號",',
				'AgentID VARCHAR(36) NOT NULL COMMENT "代理",',
				'Bet DECIMAL(20, 4) ,',
				'Win DECIMAL(20, 4) ,',
				'Content TEXT COMMENT "下注內容",',
				'CreatAt DATETIME NOT NULL'
            ')'
        );
		PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END ;
	`
	e = s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}
}
