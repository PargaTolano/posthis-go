DROP DATABASE IF EXISTS posthis_local;
CREATE DATABASE posthis_local;

USE posthis_local;

DROP PROCEDURE IF EXISTS SP_GET_FEED;
DELIMITER $$
CREATE PROCEDURE SP_GET_FEED 
(
	IN in_id int,
	IN in_offset int, 
	IN in_limit int
)
BEGIN
DROP TABLE IF EXISTS temp_FOLLOWED_IDS;
	CREATE TABLE temp_FOLLOWED_IDS(id int) 
    AS
    SELECT u.id AS ID FROM users u
    JOIN   follows f ON u.ID = f.followed_id AND f.follower_id = in_id;
    
	DROP TABLE IF EXISTS temp_REPOSTED_IDS;
	CREATE TABLE temp_REPOSTED_IDS(id int) 
    AS
    SELECT p.id
    FROM posts p
    JOIN reposts r 	ON r.post_id = p.id
    WHERE r.user_id in (SELECT id FROM temp_FOLLOWED_IDS);
    
    CREATE TEMPORARY TABLE temp_FINAL_POSTS
    AS
    SELECT 
		p.id 						ID,
		u1.id 						PublisherID,
        u1.username					PublisherUserName,
        u1.tag      				PublisherTag,
        COALESCE(m.name, '')		PublisherProfilePic,
		COALESCE(u2.username, '') 	ReposterID,
        COALESCE(u2.tag, '')		ReposterTag,
		IF( COALESCE(p.id in (SELECT id FROM temp_REPOSTED_IDS), FALSE), r.created_at, p.created_at) Date,
        p.Content 					Content,
		COALESCE(r.id, 0) 			RepostID,
        p.id in (SELECT id FROM temp_REPOSTED_IDS) 						IsRepost,
        COALESCE(u2.id in (SELECT id from temp_FOLLOWED_IDS), FALSE) 	IsFollowedRepost,
        u1.id in (SELECT id from temp_FOLLOWED_IDS) 					IsFollowedPost
    FROM posts p
    JOIN users u1 ON p.owner_id = u1.id
	LEFT JOIN media m  ON m.owner_id = u1.id
    LEFT JOIN reposts r ON p.id = r.post_id
    LEFT JOIN users u2 ON r.user_id = u2.id
    GROUP BY p.id
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc
    LIMIT in_offset, in_limit;
    
    SELECT 
		ID, 
		PublisherID, 
		PublisherUserName, 
		PublisherTag, 
		PublisherProfilePic, 
		ReposterID, 
		ReposterTag, 
		Date, 
		Content, 
		RepostID
    FROM temp_FINAL_POSTS
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc;
    
	SELECT 
		ID
    FROM temp_FINAL_POSTS
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc;
    
    DROP TABLE temp_FINAL_POSTS;
END $$
DELIMITER ;
DROP PROCEDURE IF EXISTS SP_GET_USER_FEED;
DELIMITER $$
CREATE PROCEDURE SP_GET_USER_FEED 
(
	IN in_id int, 
    IN in_poster_id int,
	IN in_offset int, 
	IN in_limit int
)
BEGIN
	DROP TABLE IF EXISTS temp_FOLLOWED_IDS;
	CREATE TABLE temp_FOLLOWED_IDS(id int) 
    AS
    SELECT u.id AS ID FROM users u
    JOIN   follows f ON u.ID = f.followed_id AND f.follower_id = in_id;
    
	DROP TABLE IF EXISTS temp_REPOSTED_IDS;
	CREATE TABLE temp_REPOSTED_IDS(id int) 
    AS
    SELECT p.id
    FROM posts p
    JOIN reposts r 	ON r.post_id = p.id
    WHERE r.user_id in (SELECT id FROM temp_FOLLOWED_IDS);
    
    CREATE TEMPORARY TABLE temp_FINAL_POSTS
    AS
    SELECT 
		p.id 						ID,
		u1.id 						PublisherID,
        u1.username					PublisherUserName,
        u1.tag      				PublisherTag,
        COALESCE(m.name, '')		PublisherProfilePic,
		COALESCE(u2.username, '') 	ReposterID,
        COALESCE(u2.tag, '')		ReposterTag,
		IF( COALESCE(p.id in (SELECT id FROM temp_REPOSTED_IDS), FALSE), r.created_at, p.created_at) Date,
        p.Content 					Content,
		COALESCE(r.id, 0) 			RepostID,
        p.id in (SELECT id FROM temp_REPOSTED_IDS) 						IsRepost,
        COALESCE(u2.id in (SELECT id from temp_FOLLOWED_IDS), FALSE) 	IsFollowedRepost,
        u1.id in (SELECT id from temp_FOLLOWED_IDS) 					IsFollowedPost
    FROM posts p
    JOIN users u1 ON p.owner_id = u1.id
	LEFT JOIN media m  ON m.owner_id = u1.id
    LEFT JOIN reposts r ON p.id = r.post_id
    LEFT JOIN users u2 ON r.user_id = u2.id
	WHERE IF( IsRepost, u2.id, u1.id) = in_poster_id
    GROUP BY p.id
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc
    LIMIT in_offset, in_limit;
    
    SELECT 
		ID, 
		PublisherID, 
		PublisherUserName, 
		PublisherTag, 
		PublisherProfilePic, 
		ReposterID, 
		ReposterTag, 
		Date, 
		Content, 
		RepostID
    FROM temp_FINAL_POSTS
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc;
    
	SELECT 
		ID
    FROM temp_FINAL_POSTS
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc;
    
    DROP TABLE temp_FINAL_POSTS;
    
END $$
DELIMITER ;

DROP PROCEDURE IF EXISTS SP_SEARCH_POSTS;
DELIMITER $$
CREATE PROCEDURE SP_SEARCH_POSTS 
(
	IN in_query    	LONGTEXT, 
    #IN in_hashtags 	LONGTEXT,
	IN in_offset 	int, 
	IN in_limit 	int
)
BEGIN
	
    CREATE TEMPORARY TABLE IF NOT EXISTS temp_SEARCH_POSTS
    AS
	SELECT 
		p.id 					ID,
        p.content 				Content,
		u.id 					PublisherID,
		u.username 				PublisherUserName,
        u.tag      				PublisherTag,
        COALESCE( m.name, '')	PublisherProfilePic,
        p.created_at			Date,
        COUNT(l.id)				LikeCount,
        COUNT(ry.id)			ReplyCount,
        COUNT(r.id)				RepostsCount
	FROM posts p
	JOIN users u ON p.owner_id = u.id
    LEFT JOIN media m 		ON u.id = m.owner_id
    LEFT JOIN likes l 		ON l.post_id  = p.id
    LEFT JOIN replies ry 	ON ry.post_id = p.id
    LEFT JOIN reposts r		ON r.post_id  = p.id
    WHERE p.content LIKE CONCAT('%', in_query,'%')
    GROUP BY ID, Content, PublisherID, PublisherUserName, PublisherTag, PublisherProfilePic;
    
    SELECT * FROM temp_SEARCH_POSTS
    LIMIT in_offset, in_limit;
    
    SELECT ID FROM temp_SEARCH_POSTS
    LIMIT in_offset, in_limit;
    
END $$
DELIMITER ;

DROP PROCEDURE IF EXISTS SP_SEARCH_USERS;
DELIMITER $$
CREATE PROCEDURE SP_SEARCH_USERS 
(
	IN in_query    	LONGTEXT, 
    #IN in_hashtags 	LONGTEXT,
	IN in_offset 	int, 
	IN in_limit 	int
)
BEGIN
	
    CREATE TEMPORARY TABLE IF NOT EXISTS temp_SEARCH_USERS
    AS
	SELECT 
		u.id 					ID,
        u.tag 					Tag,
        u.username 				Username,
        COALESCE( m.name, '') 	ProfilePicPath
	FROM users u
    LEFT JOIN media m ON u.id = m.owner_id
    WHERE u.username LIKE CONCAT('%', in_query,'%') OR  u.tag LIKE CONCAT('%', in_query,'%');
    
    SELECT * FROM temp_SEARCH_USERS
    LIMIT in_offset, in_limit;
    
    DROP  TEMPORARY TABLE temp_SEARCH_USERS;
    
END $$
DELIMITER ;