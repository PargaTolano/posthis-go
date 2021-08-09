DROP DATABASE IF EXISTS posthis_local;
CREATE DATABASE posthis_local;

USE posthis_local;

DROP PROCEDURE IF EXISTS SP_GET_PROFILE;
DELIMITER $$
CREATE PROCEDURE SP_GET_PROFILE
(
	IN in_id 			INT,
    IN in_viewer_id 	INT
)
BEGIN
	SELECT 
		u.id 					id,
		u.username				username,
        u.tag					tag,
        u.email					email,
        COALESCE(m1.name,'') 	profilepicpath,
        COALESCE(m2.name,'')	coverpicpath,
        COUNT(f1.id) 			followerCount,
        COUNT(f2.id)			followingCount,
		f3.id IS NOT NULL 		isFollowed
	FROM users u
    LEFT JOIN follows f1 ON f1.followed_id = in_id
    LEFT JOIN follows f2 ON f1.follower_id = in_id
    LEFT JOIN media m1 ON u.id = m1.owner_id AND m1.owner_type = 'profilepicuser'
    LEFT JOIN media m2 ON u.id = m2.owner_id AND m2.owner_type = 'coverpicuser'
    LEFT JOIN follows f3 ON f3.followed_id = in_id AND f3.follower_id = in_viewer_id
	WHERE u.id = in_id
    GROUP BY u.id;
END $$
DELIMITER ;
DROP PROCEDURE IF EXISTS SP_GET_POST_DETAIL;
DELIMITER $$
CREATE PROCEDURE SP_GET_POST_DETAIL
(
	IN in_id int,
    IN in_viewer_id int
)
BEGIN
	
    #MAKE SELECTION FOLLOW THE SAME FORMAT AS FEED POSTS BUT NO REPOST INFO
	SELECT 
		p.id 												ID,
		u1.id 												PublisherID,
        u1.username											PublisherUserName,
        u1.tag      										PublisherTag,
        COALESCE(m.name, '')								PublisherProfilePic,
		'' 													ReposterID,
        ''													ReposterTag,
		p.created_at 										Date,
        p.Content 											Content,
		0 													RepostID,
        EXISTS ( SELECT*FROM likes   WHERE post_id = p.id AND user_id = in_viewer_id AND deleted_at IS NULL )   IsLiked,
        EXISTS ( SELECT*FROM reposts WHERE post_id = p.id AND user_id = in_viewer_id AND deleted_at IS NULL )   IsReposted
    FROM posts p
    JOIN users u1 ON p.owner_id = u1.id
	LEFT JOIN media m  ON m.owner_id = u1.id AND m.owner_type='profilepicuser'
    WHERE p.id = in_id
    GROUP BY ID;
	
END $$
DELIMITER ;

DROP PROCEDURE IF EXISTS SP_GET_POST_REPLIES;
DELIMITER $$
CREATE PROCEDURE SP_GET_POST_REPLIES
(
	IN in_id int
)
BEGIN
	
    CREATE TEMPORARY TABLE temp_t
    AS
	SELECT 
		ry.id 												ReplyID,
        ry.Content 											Content,
        ry.post_id											PostID,
		u1.id 												PublisherID,
        u1.username											PublisherUserName,
        u1.tag      										PublisherTag,
        COALESCE(m.name, '')								PublisherProfilePic,
		ry.created_at 										Date
    FROM replies ry
    JOIN users u1 ON ry.user_id = u1.id
	LEFT JOIN media m  ON m.owner_id = u1.id AND m.owner_type='profilepicuser'
    WHERE ry.post_id = in_id
    GROUP BY ReplyID
    ORDER BY Date DESC;
    
    SELECT * FROM temp_t
    ORDER BY Date DESC;
    
    SELECT ReplyID FROM temp_t
    ORDER BY Date DESC;
    
    DROP TEMPORARY TABLE temp_t;
	
END $$
DELIMITER ;

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
	CREATE TABLE temp_FOLLOWED_IDS
    AS
    SELECT u.id AS ID FROM users u
    JOIN follows f ON u.ID = f.followed_id AND f.follower_id = in_id;
    
	DROP TABLE IF EXISTS temp_REPOSTED_IDS;
	CREATE TABLE temp_REPOSTED_IDS 
    AS
    SELECT p.id
    FROM posts p
    JOIN reposts r 	ON r.post_id = p.id
    WHERE r.user_id in (SELECT id FROM temp_FOLLOWED_IDS);
    
    CREATE TEMPORARY TABLE temp_FINAL_POSTS
    AS
    SELECT 
		p.id 																										ID,
		u1.id 																										PublisherID,
        u1.username																									PublisherUserName,
        u1.tag      																								PublisherTag,
        COALESCE(m.name, '')																						PublisherProfilePic,
		COALESCE(u2.username, '') 																					ReposterID,
        COALESCE(u2.tag, '')																						ReposterTag,
		IF( p.id IN (SELECT id FROM temp_REPOSTED_IDS), r.created_at, p.created_at) 			 					Date,
        p.Content 																								 	Content,
		COALESCE(r.id, 0) 																						 	RepostID,
        p.id in (SELECT id FROM temp_REPOSTED_IDS) 																 	IsRepost,
        COALESCE(u2.id in (SELECT id from temp_FOLLOWED_IDS), FALSE) 											 	IsFollowedRepost,
        u1.id in (SELECT id from temp_FOLLOWED_IDS) 															 	IsFollowedPost,
        EXISTS ( SELECT*FROM likes l 	WHERE l.post_id = p.id AND l.user_id = in_id AND l.deleted_at IS NULL )	 	IsLiked,
        EXISTS ( SELECT*FROM reposts r 	WHERE post_id   = p.id AND r.user_id = in_id AND r.deleted_at IS NULL )  	IsReposted
    FROM posts p
    JOIN users u1 ON p.owner_id = u1.id AND u1.deleted_at IS NULL
	LEFT JOIN media m  ON m.owner_id = u1.id AND m.owner_type = 'profilepicuser'
    LEFT JOIN reposts r ON p.id = r.post_id AND r.user_id != in_id
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
		RepostID,
        IsLiked,
        IsReposted
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
		p.id 																										ID,
		u1.id 																										PublisherID,
        u1.username																									PublisherUserName,
        u1.tag      																								PublisherTag,
        COALESCE(m.name, '')																						PublisherProfilePic,
		COALESCE(u2.username, '') 																					ReposterID,
        COALESCE(u2.tag, '')																						ReposterTag,
		IF( p.id IN (SELECT id FROM temp_REPOSTED_IDS), r.created_at, p.created_at) 			 					Date,
        p.Content 																								 	Content,
		COALESCE(r.id, 0) 																						 	RepostID,
        p.id in (SELECT id FROM temp_REPOSTED_IDS) 																 	IsRepost,
        COALESCE(u2.id in (SELECT id from temp_FOLLOWED_IDS), FALSE) 											 	IsFollowedRepost,
        u1.id in (SELECT id from temp_FOLLOWED_IDS) 															 	IsFollowedPost,
        EXISTS ( SELECT*FROM likes l 	WHERE l.post_id = p.id AND l.user_id = in_id AND l.deleted_at IS NULL )	 	IsLiked,
        EXISTS ( SELECT*FROM reposts r 	WHERE post_id   = p.id AND r.user_id = in_id AND r.deleted_at IS NULL )  	IsReposted
    FROM posts p
    JOIN users u1 ON p.owner_id = u1.id AND u1.deleted_at IS NULL
	LEFT JOIN media m  ON m.owner_id = u1.id AND m.owner_type = 'profilepicuser'
    LEFT JOIN reposts r ON p.id = r.post_id AND r.user_id != in_id
    LEFT JOIN users u2 ON r.user_id = u2.id
	WHERE IF( p.id in (SELECT id FROM temp_REPOSTED_IDS) , u2.id, u1.id) = in_poster_id
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
		RepostID,
        IsLiked,
        IsReposted
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
		p.id 												ID,
        p.content 											Content,
		u.id 												PublisherID,
		u.username 											PublisherUserName,
        u.tag      											PublisherTag,
        COALESCE( m.name, '')								PublisherProfilePic,
        p.created_at										Date,
        COUNT(l.id)											LikeCount,
        COUNT(ry.id)										ReplyCount,
        COUNT(r.id)											RepostsCount,
		EXISTS ( SELECT*FROM likes WHERE post_id = p.id AND deleted_at IS NULL ) 	IsLiked,
        EXISTS ( SELECT*FROM reposts WHERE post_id = p.id AND deleted_at IS NULL )   IsReposted
	FROM posts p
	JOIN users u 			ON p.owner_id = u.id
    LEFT JOIN media m 		ON u.id = m.owner_id AND m.owner_type = 'profilepicuser'
    LEFT JOIN likes l 		ON l.post_id  = p.id
    LEFT JOIN replies ry 	ON ry.post_id = p.id
    LEFT JOIN reposts r		ON r.post_id  = p.ids
    WHERE p.content LIKE CONCAT('%', in_query,'%')
    GROUP BY ID, PublisherUserName, PublisherTag, PublisherProfilePic
    ORDER BY Date
    LIMIT in_offset, in_limit;
    
    SELECT * FROM temp_SEARCH_POSTS
    ORDER BY Date;
    
    SELECT ID FROM temp_SEARCH_POSTS
    ORDER BY Date;
    
    DROP  TEMPORARY TABLE temp_SEARCH_POSTS;
    
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
        u.username 				UserName,
        COALESCE( m.name, '') 	ProfilePicPath
	FROM users u
    LEFT JOIN media m ON u.id = m.owner_id
    WHERE u.username LIKE CONCAT('%', in_query,'%') OR  u.tag LIKE CONCAT('%', in_query,'%')
    GROUP BY ID
    ORDER BY UserName
    LIMIT in_offset, in_limit;
    
    SELECT 
		ID,
        Tag,
        Username,
        ProfilePicPath
    FROM temp_SEARCH_USERS
    ORDER BY UserName;
    
    DROP  TEMPORARY TABLE temp_SEARCH_USERS;
    
END $$
DELIMITER ;

#RUN THIS PART AFTER INITIALIZING THE GO APP
CREATE INDEX idx_post_id ON posts(id);

SELECT * FROM POSTS;

#TEST
DROP PROCEDURE IF EXISTS TSP_TEST_FEED;
DELIMITER $$
CREATE PROCEDURE TSP_TEST_FEED
(IN IN_ID INT, IN IN_OFFSET INT, IN IN_LIMIT INT)
BEGIN

	DROP TABLE IF EXISTS temp_FOLLOWED_IDS;
	CREATE TABLE temp_FOLLOWED_IDS(id int) 
    AS
    SELECT u.id AS ID FROM users u
    JOIN   follows f ON u.ID = f.followed_id AND f.follower_id = IN_ID AND f.deleted_at IS NOT NULL;
    
	DROP TABLE IF EXISTS temp_REPOSTED_IDS;
	CREATE TABLE temp_REPOSTED_IDS(id int) 
    AS
    SELECT p.id
    FROM posts p
    JOIN reposts r 	ON r.post_id = p.id
    WHERE r.user_id in (SELECT id FROM temp_FOLLOWED_IDS);
    
    #SELECT * FROM temp_FOLLOWED_IDS;
    #SELECT * FROM temp_REPOSTED_IDS;

	SELECT 
		p.id 																										ID,
		u1.id 																										PublisherID,
        u1.username																									PublisherUserName,
        u1.tag      																								PublisherTag,
        COALESCE(m.name, '')																						PublisherProfilePic,
		COALESCE(u2.username, '') 																					ReposterID,
        COALESCE(u2.tag, '')																						ReposterTag,
		IF( p.id IN (SELECT id FROM temp_REPOSTED_IDS), r.created_at, p.created_at) 			 					Date,
        p.Content 																								 	Content,
		COALESCE(r.id, 0) 																						 	RepostID,
        p.id in (SELECT id FROM temp_REPOSTED_IDS) 																 	IsRepost,
        COALESCE(u2.id in (SELECT id from temp_FOLLOWED_IDS), FALSE) 											 	IsFollowedRepost,
        u1.id in (SELECT id from temp_FOLLOWED_IDS) 															 	IsFollowedPost,
        EXISTS ( SELECT*FROM likes l 	WHERE l.post_id = p.id AND l.user_id = in_id AND l.deleted_at IS NULL )	 	IsLiked,
        EXISTS ( SELECT*FROM reposts r 	WHERE post_id   = p.id AND r.user_id = in_id AND r.deleted_at IS NULL )  	IsReposted
    FROM posts p
    JOIN users u1 ON p.owner_id = u1.id AND u1.deleted_at IS NULL
	LEFT JOIN media m  ON m.owner_id = u1.id AND m.owner_type = 'profilepicuser'
    LEFT JOIN reposts r ON p.id = r.post_id AND r.user_id != in_id
    LEFT JOIN users u2 ON r.user_id = u2.id
    GROUP BY p.id
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc
    LIMIT in_offset, in_limit;
    
END $$
DELIMITER ;

CALL TSP_TEST_FEED(1,0,1);

SELECT users.* FROM posts LEFT OUTER JOIN users ON 1=1;