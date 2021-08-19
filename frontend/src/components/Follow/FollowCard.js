import React,{ useState } from 'react';

import {
    Button,
}from  '@material-ui/core';


import { createFollow, deleteFollow } from '_api';

import styles from '_styles/FollowCard.module.css';

import defaultProfilePic  from 'assets/avatar-placeholder.svg';
import { Link } from 'react-router-dom';
import { routes } from '_utils';
import { authenticationService } from '_services';

export const FollowCard = ({user:userProp}) => {

    const [user, setUser] = useState(userProp);

    const onClickFollowButton = async ()=>{
        const {data:responseData, err} = user.isFollowed ? await deleteFollow( user.id ) : await createFollow( user.id );

        if ( err !== null ) return;

        const {data} = responseData;
        setUser(data);
    };

    return (
        <div className={styles.root}>
            <Link
                className={styles.profileLink}
                to={routes.getProfile(user.id)}
            >
                <img className={styles.pfp} src={user.profilePicPath || defaultProfilePic } />
            </Link>

            <Link
                className={styles.profileLink} 
                to={routes.getProfile(user.id)}
            >
                <div className={styles.detail}> 
                    <h4 className={styles.userdata}>
                        {user.username}@<span className={styles.tag}>{user.tag}</span>
                    </h4>
                </div>
            </Link>
            
            {   
                authenticationService.currentUserValue.id !== user.id &&
                <Button
                    variant='contained'
                    color='secondary'
                    className={styles.followBtn}
                    onClick={onClickFollowButton}
                >
                    { user.isFollowed ? 'Unfollow' : 'Follow' }
                </Button>
            }
        </div>
    )
}

export default FollowCard;
