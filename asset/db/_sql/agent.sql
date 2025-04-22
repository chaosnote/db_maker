/*
建立代理商表單
*/
CREATE TABLE IF NOT EXISTS agent (
  ID VARCHAR(36) NOT NULL,
  Level INT UNSIGNED, /* {0:預設,...} */
  Name VARCHAR(20) NOT NULL,
  APIKey VARCHAR(8) NOT NULL,
  Category VARCHAR(10) NOT NULL, /* 開放的遊戲分類 */
  ThirdParty TEXT, /* 對方提供的資訊 */
  PRIMARY KEY (ID),
  UNIQUE KEY Name (Name)
);