import React from 'react';
import { Link } from 'react-router-dom';

import { followService}      from '_services';
import { routes } from '_utils';

import defaultImage from 'assets/avatar-placeholder.svg';
import styles from '_styles/UserCard.module.css';
import { prettyMagnitude } from '_helpers';

const ChannelPhoto = ({ user }) => {
    return(
        <Link
            to={routes.getProfile(user.id)}
        >
            <img className = {styles.pfp} src ={ user.profilePicPath || defaultImage}/>
        </Link>
    );
};

const ChannelFollows = ({ user })=>{
    return(
        <div className={styles.followContainer}>
            <div className={styles.followCount}>
                <span 
                    className={styles.followLink} 
                    onClick   = { ()=>followService.getFollowerUsers(user.id) }
                >
                    { prettyMagnitude(user.followerCount) } followers
                </span>
            </div>
            <div className={styles.followCount}>
                <span 
                    className={styles.followLink} 
                    onClick   = { ()=>followService.getFollowedUsers(user.id) }
                >
                    { prettyMagnitude(user.followingCount) } following   
                </span>     
            </div>
        </div>
    );
};

const ChannelUsername = ({ user }) => {
    return(
        <div className={styles.userinfo}>
            <Link 
                className={styles.infoLink} 
                to={routes.getProfile(user.id)}
            >
                <p className={styles.username}>
                    {user.username}
                </p>
                <p className={styles.tag}>
                    {`@${user.tag}`}
                </p>
            </Link>
            <ChannelFollows user={user} />
        </div>
    );
};

const SeeButton = ({ user }) => {
    return(
        <Link 
            className={styles.profileLink} 
            to={routes.getProfile(user.id)}
        >
            <button className={styles.btn}>
                GO TO PROFILE
            </button>
        </Link>
    );
};

export const SearchUserCard = ( {user} ) => {
    
    return(
        <div className={styles.root}>
            <ChannelPhoto       user={user}/>
            <ChannelUsername    user={user}/>
            <SeeButton          user={user}/>
        </div>
    );
}

export default SearchUserCard;