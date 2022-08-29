SET SESSION sql_mode=(SELECT REPLACE(@@sql_mode,'ONLY_FULL_GROUP_BY',''));

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
        COUNT(DISTINCT f1.id) 			followerCount,
        COUNT(DISTINCT f2.id)			followingCount,
		f3.id IS NOT NULL 		isFollowed
	FROM users u
    LEFT JOIN follows f1 ON f1.followed_id = in_id
    LEFT JOIN follows f2 ON f2.follower_id = in_id
    LEFT JOIN follows f3 ON f3.followed_id = in_id AND f3.follower_id = in_viewer_id
    LEFT JOIN media m1 ON u.id = m1.owner_id AND m1.owner_type = 'profilepicuser'
    LEFT JOIN media m2 ON u.id = m2.owner_id AND m2.owner_type = 'coverpicuser'
	WHERE u.id = in_id and u.deleted_at IS NULL
    GROUP BY u.id;
END $$
DELIMITER ;

DROP PROCEDURE IF EXISTS SP_GET_FOLLOWERS;
DELIMITER $$
CREATE PROCEDURE SP_GET_FOLLOWERS
(
	IN in_id 			INT,
    IN in_viewer_id 	INT
)
BEGIN
	SELECT 
		u.id 												id,
		u.username											username,
        u.tag												tag,
        COALESCE(m.name,'') 								profilepicpath,
		MAX( COALESCE( f2.id <> in_viewer_id, FALSE )) > 0  isFollowed
	FROM 		users u
    JOIN		follows f 	ON f.follower_id  = u.id AND f.followed_id = in_id
    LEFT JOIN 	media	m 	ON m.owner_id = u.id AND m.owner_type = 'profilepicuser'
    LEFT JOIN 	follows f2 	ON f2.follower_id = in_viewer_id
    GROUP BY u.id;
END $$
DELIMITER ;

DROP PROCEDURE IF EXISTS SP_GET_FOLLOWING;
DELIMITER $$
CREATE PROCEDURE SP_GET_FOLLOWING
(
	IN in_id 			INT,
    IN in_viewer_id 	INT
)
BEGIN
	SELECT 
		u.id 								id,
		u.username							username,
        u.tag								tag,
        COALESCE(m.name,'') 				profilepicpath,
		MAX( COALESCE( f2.id <> in_viewer_id, FALSE )) > 0  isFollowed
	FROM 		users u
    JOIN		follows f 	ON f.followed_id  = u.id AND f.follower_id = in_id
    LEFT JOIN 	media	m 	ON m.owner_id = u.id AND m.owner_type = 'profilepicuser'
    LEFT JOIN 	follows f2 	ON f2.follower_id = in_viewer_id
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
		p.id 								ID,									
		u1.id 								PublisherID,
        u1.username							PublisherUserName,
        u1.tag      						PublisherTag,
        COALESCE(m.name, '')				PublisherProfilePic,
		'' 									ReposterUsername,
        ''									ReposterTag,
		p.created_at 						Date,
        p.Content 							Content,
		0 									RepostID,
        COUNT(DISTINCT l.id)				LikeCount,
        COUNT(DISTINCT ry.id)				ReplyCount,
        COUNT(DISTINCT r.id)				RepostsCount,
		MAX( COALESCE( l2.id, FALSE )) > 0  IsLiked,
		MAX( COALESCE( r2.id, FALSE )) > 0	IsReposted
    FROM posts p
    JOIN users u1 ON p.owner_id = u1.id AND u1.deleted_at IS NULL
	LEFT JOIN media    m  ON m.owner_id = u1.id AND m.owner_type='profilepicuser'
    LEFT JOIN likes    l  ON  l.post_id = p.id
    LEFT JOIN replies ry  ON ry.post_id = p.id
    LEFT JOIN reposts  r  ON  r.post_id = p.id
    LEFT JOIN likes   l2  ON p.id  = l2.post_id AND l2.user_id = in_viewer_id
    LEFT JOIN reposts r2  ON p.id  = r2.post_id AND r2.user_id = in_viewer_id
    WHERE p.id = in_id AND p.deleted_at IS NULL
    GROUP BY p.id;
	
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
    JOIN users u1 ON ry.user_id = u1.id AND u1.deleted_at IS NULL
	LEFT JOIN media m  ON m.owner_id = u1.id AND m.owner_type='profilepicuser'
    WHERE ry.post_id = in_id
    GROUP BY ReplyID
    ORDER BY Date;
    
    SELECT * FROM temp_t
    ORDER BY Date;
    
    SELECT ReplyID FROM temp_t
    ORDER BY Date;
    
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
    DROP TEMPORARY TABLE IF EXISTS temp_FINAL_POSTS;
	DROP TEMPORARY TABLE IF EXISTS temp_REPOSTED_IDS;
	DROP TEMPORARY TABLE IF EXISTS temp_FOLLOWED_IDS;
    
	CREATE 	TEMPORARY TABLE temp_FOLLOWED_IDS(id int)
    AS
    SELECT u.id ID 
    FROM users u
    JOIN follows f ON u.ID = f.followed_id AND f.follower_id = in_id;
    
	CREATE 	TEMPORARY TABLE temp_REPOSTED_IDS(id int) 
    AS
    SELECT 	p.id
    FROM 	posts p
    JOIN 	users u 				ON u.id = p.owner_id
    JOIN 	reposts r 				ON r.post_id = p.id
    JOIN  	temp_FOLLOWED_IDS fids 	ON fids.id = u.id
    GROUP BY p.id;
    
    CREATE TEMPORARY TABLE temp_FINAL_POSTS
    AS
    SELECT 
		p.id 																										ID,
		u1.id 																										PublisherID,
        u1.username																									PublisherUserName,
        u1.tag      																								PublisherTag,
        COALESCE(m.name, '')																						PublisherProfilePic,
		COALESCE(u2.username, '') 																					ReposterUsername,
        COALESCE(u2.tag, '')																						ReposterTag,
		IF( p.id = rids.id, r.created_at, p.created_at) 			 												Date,
        p.Content 																								 	Content,
		COALESCE(r.id, 0) 																						 	RepostID,
		COALESCE(p.id  = rids.id, FALSE)																			IsRepost,
        COALESCE(u2.id = fids.id, FALSE) 											 								IsFollowedRepost,
        u1.id = fids.id 																						 	IsFollowedPost,
        COUNT(DISTINCT l2.id)																						LikeCount,
        COUNT(DISTINCT ry.id)																						ReplyCount,
        COUNT(DISTINCT r3.id)																						RepostsCount,
        MAX( COALESCE(  l.id, FALSE)) > 0																			IsLiked,
        MAX( COALESCE( r2.id, FALSE)) > 0																  			IsReposted
    FROM posts p
    JOIN users u1 						ON p.owner_id 	= u1.id 	 AND u1.deleted_at 	IS NULL
    LEFT JOIN likes   l 				ON l.post_id 	= p.id		 AND l.user_id = in_id
    LEFT JOIN likes   l2 				ON l2.post_id 	= p.id
    LEFT JOIN replies ry 				ON ry.post_id 	= p.id
	LEFT JOIN media   m  				ON m.owner_id 	= u1.id 	 AND m.owner_type = 'profilepicuser'
    LEFT JOIN reposts r 				ON p.id 		= r.post_id  AND r.user_id != in_id
    LEFT JOIN reposts r2 				ON p.id 		= r2.post_id AND r2.user_id = in_id
    LEFT JOIN reposts r3 				ON p.id 		= r3.post_id
    LEFT JOIN users   u2 				ON r.user_id 	= u2.id
    LEFT JOIN temp_REPOSTED_IDS rids 	ON rids.id 		= p.id
    LEFT JOIN temp_FOLLOWED_IDS fids 	ON fids.id 		= u1.id OR fids.id = u2.id
    WHERE p.deleted_at IS NULL
    GROUP BY p.id
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc
    LIMIT in_offset, in_limit;
    
    SELECT
		ID,
		PublisherID,
		PublisherUserName,
		PublisherTag,
		PublisherProfilePic,
		ReposterUsername,
		ReposterTag,
		Date,
		Content,
		RepostID,
        LikeCount,
        ReplyCount,
        RepostsCount,
        IsLiked,
        IsReposted
    FROM temp_FINAL_POSTS
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc;
    
	SELECT 
		ID
    FROM temp_FINAL_POSTS
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc;
    
    DROP TEMPORARY TABLE temp_FINAL_POSTS;
	DROP TEMPORARY TABLE temp_REPOSTED_IDS;
	DROP TEMPORARY TABLE temp_FOLLOWED_IDS;
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
    DROP TEMPORARY TABLE IF EXISTS temp_FINAL_POSTS;
	DROP TEMPORARY TABLE IF EXISTS temp_REPOSTED_IDS;
	DROP TEMPORARY TABLE IF EXISTS temp_FOLLOWED_IDS;
    
	CREATE 	TEMPORARY TABLE temp_FOLLOWED_IDS(id int)
    AS
    SELECT u.id ID 
    FROM users u
    JOIN follows f ON u.ID = f.followed_id AND f.follower_id = in_id;
    
	CREATE 	TEMPORARY TABLE temp_REPOSTED_IDS(id int) 
    AS
    SELECT 	p.id
    FROM 	posts p
    JOIN 	users u 				ON u.id = p.owner_id
    JOIN 	reposts r 				ON r.post_id = p.id
    JOIN  	temp_FOLLOWED_IDS fids 	ON fids.id = u.id
    GROUP BY p.id;
    
    CREATE TEMPORARY TABLE temp_FINAL_POSTS
    AS
    SELECT 
		p.id 																										ID,
		u1.id 																										PublisherID,
        u1.username																									PublisherUserName,
        u1.tag      																								PublisherTag,
        COALESCE(m.name, '')																						PublisherProfilePic,
		COALESCE(u2.username, '') 																					ReposterUsername,
        COALESCE(u2.tag, '')																						ReposterTag,
		IF( p.id = rids.id, r.created_at, p.created_at) 			 												Date,
        p.Content 																								 	Content,
		COALESCE(r.id, 0) 																						 	RepostID,
		COALESCE(p.id  = rids.id, FALSE)																			IsRepost,
        COALESCE(u2.id = fids.id, FALSE) 											 								IsFollowedRepost,
        u1.id = fids.id 																						 	IsFollowedPost,
        COUNT(DISTINCT l2.id)																						LikeCount,
        COUNT(DISTINCT ry.id)																						ReplyCount,
        COUNT(DISTINCT r3.id)																						RepostsCount,
        MAX( COALESCE(  l.id, FALSE)) > 0																			IsLiked,
        MAX( COALESCE( r2.id, FALSE)) > 0																  			IsReposted
    FROM posts p
    JOIN users u1 						ON p.owner_id 	= u1.id 	 AND u1.deleted_at 	IS NULL
    LEFT JOIN likes   l 				ON l.post_id 	= p.id		 AND l.user_id = in_id
    LEFT JOIN likes   l2 				ON l2.post_id 	= p.id
    LEFT JOIN replies ry 				ON ry.post_id 	= p.id
	LEFT JOIN media   m  				ON m.owner_id 	= u1.id 	 AND m.owner_type = 'profilepicuser'
    LEFT JOIN reposts r 				ON p.id 		= r.post_id  AND r.user_id != in_id
    LEFT JOIN reposts r2 				ON p.id 		= r2.post_id AND r2.user_id = in_id
    LEFT JOIN reposts r3 				ON p.id 		= r3.post_id
    LEFT JOIN users   u2 				ON r.user_id 	= u2.id
    LEFT JOIN temp_REPOSTED_IDS rids 	ON rids.id 		= p.id
    LEFT JOIN temp_FOLLOWED_IDS fids 	ON fids.id 		= u1.id OR fids.id = u2.id
    WHERE IF( p.id = rids.id , u2.id, u1.id) = in_poster_id  AND p.deleted_at IS NULL
    GROUP BY p.id
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc
    LIMIT in_offset, in_limit;
    
    SELECT 
		ID, 
		PublisherID, 
		PublisherUserName, 
		PublisherTag, 
		PublisherProfilePic, 
		ReposterUsername, 
		ReposterTag, 
		Date, 
		Content, 
		RepostID,
        LikeCount,
        ReplyCount,
        RepostsCount,
        IsLiked,
        IsReposted
    FROM temp_FINAL_POSTS
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc;
    
	SELECT 
		ID
    FROM temp_FINAL_POSTS
    ORDER BY Date DESC, IsFollowedRepost desc, IsFollowedPost desc;
    
    DROP TEMPORARY TABLE temp_FINAL_POSTS;
	DROP TEMPORARY TABLE temp_REPOSTED_IDS;
	DROP TEMPORARY TABLE temp_FOLLOWED_IDS;
    
END $$
DELIMITER ;

DROP PROCEDURE IF EXISTS SP_SEARCH_POSTS;
DELIMITER $$
CREATE PROCEDURE SP_SEARCH_POSTS 
(
	IN in_query    		LONGTEXT, 
    IN in_viewer_id 	int,
	IN in_offset 		int, 
	IN in_limit 		int
)
BEGIN
	
    CREATE TEMPORARY TABLE IF NOT EXISTS temp_SEARCH_POSTS
    AS
	SELECT 
		p.id 																											ID,
        p.content 																										Content,
		u.id 																											PublisherID,
		u.username 																										PublisherUserName,
        u.tag      																										PublisherTag,
        COALESCE( m.name, '')																							PublisherProfilePic,
        p.created_at																									Date,
        COUNT(DISTINCT l.id)																							LikeCount,
        COUNT(DISTINCT ry.id)																							ReplyCount,
        COUNT(DISTINCT r.id)																							RepostsCount,
        MAX( COALESCE( l2.id, FALSE)) > 0	IsLiked,
        MAX( COALESCE( r2.id, FALSE)) > 0  IsReposted
	FROM posts p
	JOIN users u 			ON p.owner_id = u.id
    LEFT JOIN media m 		ON u.id = m.owner_id AND m.owner_type = 'profilepicuser'
    LEFT JOIN likes l 		ON l.post_id  = p.id
    LEFT JOIN likes l2 		ON l2.post_id = p.id AND l2.user_id = in_viewer_id
    LEFT JOIN replies ry 	ON ry.post_id = p.id
    LEFT JOIN reposts r		ON r.post_id  = p.id
    LEFT JOIN reposts r2	ON r2.post_id = p.id AND r2.user_id = in_viewer_id
    WHERE p.content LIKE CONCAT('%', in_query,'%') AND p.deleted_at IS NULL
    GROUP BY p.id
    ORDER BY Date
    LIMIT in_offset, in_limit;
    
    SELECT *
	FROM temp_SEARCH_POSTS
    ORDER BY Date;
    
    SELECT id
	FROM temp_SEARCH_POSTS
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
		u.id 							ID,
        u.tag 							Tag,
        u.username 						UserName,
        COALESCE( m.name, '') 			ProfilePicPath,
		COUNT(DISTINCT f1.id)			FollowerCount,
		COUNT(DISTINCT f2.id)			FollowingCount
	FROM users u
    LEFT JOIN media m ON u.id = m.owner_id AND m.owner_type = 'profilepicuser'
    LEFT JOIN follows f1 ON f1.followed_id = u.id
    LEFT JOIN follows f2 ON f2.follower_id = u.id
    WHERE ( u.username LIKE CONCAT('%', in_query,'%') OR  u.tag LIKE CONCAT('%', in_query,'%') ) AND u.deleted_at IS NULL
    GROUP BY ID
    ORDER BY UserName
    LIMIT in_offset, in_limit;
    
    SELECT 
		ID,
        Tag,
        Username,
        ProfilePicPath,
        FollowerCount,
        FollowingCount
    FROM temp_SEARCH_USERS
    ORDER BY UserName;
    
    DROP  TEMPORARY TABLE temp_SEARCH_USERS;
    
END $$
DELIMITER ;