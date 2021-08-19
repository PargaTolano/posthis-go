import React, { useState } from 'react';
import { Link } from 'react-router-dom';

import {
    IconButton,
} from '@material-ui/core';

import {
    Favorite as FavoriteIcon,
    QuestionAnswer as QuestionAnswerIcon,
    ReplyAll       as ReplyAllIcon,
} from '@material-ui/icons';

import { routes } from '_utils';

import defaultImage from 'assets/avatar-placeholder.svg';

import styles from '_styles/SearchPostCard.module.css';

import { prettyDate, prettyMagnitude } from '_helpers';
import { authenticationService, toastService } from '_services';
import { createLike, createRepost, deleteLike, deleteRepost } from '_api';

const ChannelUserdata = ({ post }) => {
    return(
        <div className={styles.userInfo}>
            <Link
                className={styles.userLink} 
                to={routes.getProfile(post.publisherID)}
            >
              <img src={post.publisherProfilePic || defaultImage} className={styles.avatar}/>
            </Link>
            <div>
                <Link
                    className={styles.userLink}
                    to={routes.getProfile(post.publisherID)}
                >
                    <h6 className={styles.publisher}>
                        { post.publisherUserName }<span className={styles.tag}>{`@${post.publisherTag}`}</span>
                    </h6>
                </Link>
                <Link
                    className={styles.userLink}
                    to={routes.getProfile(post.publisherID)}
                >
                    <p className={styles.date}> { prettyDate(post.date) } </p>
              
                </Link>
            </div>
        </div>
    );
};

const ChanelPostContent = ({ post }) => {
    return (
        <div className={styles.content}>
            { post.content }
        </div>
    );
};

const ChannelMediaPreview = ({post})=>{

    if (post.medias === null || post.medias.length === 0)
        return null;

    return (
        <div className={styles.mediaPreview}>
            <img className={styles.imgPreview} src={post.medias[0].path}/>
        </div>
    );
};

const ChannelInteractions = ({post, setPost})=>{

    const onClickLike = async ()=>{
        if( post.isLiked ){
    
          const {data:responseData, err}  = await deleteLike(post.postID);
    
          if ( err != null ) return;
    
          const { data } = responseData;

          if (data !== null)
          setPost(data)
        }
        else{
    
            const {data:responseData, err}  = await createLike(authenticationService.currentUserValue.id, post.postID);
    
            if ( err != null ) {
                toastService.makeToast("Couldn't like post, try again later", "error");
                return;
            };
    
            const { data } = responseData;

            if (data !== null)
                setPost(data)
        }
    };
    
    const onClickRepost = async ()=>{
    
        if( post.isReposted ) {
            const {data:responseData, err} = await deleteRepost(post.postID);
            
            if ( err != null ) {
                toastService.makeToast("Couldn't like post, try again later", "error");
                return;
            };
            
            const { data } = responseData;
            
            if (data !== null)
                setPost(data)
        }
        else{
    
            const {data:responseData, err} = await createRepost(authenticationService.currentUserValue.id, post.postID);
            
            if ( err !== null)  return;
            
            const { data } = responseData;

            if (data !== null)
                setPost(data)
        }
    };

    return (
        <div className={styles.interactions}>
            <IconButton
                className={styles.iconContainer}
                onClick={onClickLike}
            >
                <FavoriteIcon className={ post.isLiked ? styles.likeIcon : styles.grayIcon }/>
                <div className={styles.count}>
                    {prettyMagnitude(post.likeCount)}
                </div>
            </IconButton>
            <Link
                to={routes.getPost(post.postID)}
            >
                <IconButton
                    className={styles.iconContainer} 
                >
                    <QuestionAnswerIcon className={styles.commentIcon}/>
                    <div className={styles.count}>
                        {prettyMagnitude(post.replyCount)}
                    </div>
                </IconButton>
            </Link>
            <IconButton
                className={styles.iconContainer}
                onClick={onClickRepost}
            >
                <ReplyAllIcon className={ post.isReposted ? styles.repostIcon : styles.grayIcon }/>
                <div className={styles.count}>
                    {prettyMagnitude(post.repostCount)}
                </div>
            </IconButton>
        </div>
    );
};

export const SearchPostCard = ( { post:PostProp } ) => {

    const [post, setPost] = useState(PostProp);
    
    return(
        <div className={styles.root}>
            <div>
                <ChannelUserdata        post={post} />
                <ChanelPostContent      post={post} />
            </div>
            <ChannelMediaPreview    post={post} />
            <ChannelInteractions    post={post} setPost={setPost} />
        </div>
    );
}

export default SearchPostCard;