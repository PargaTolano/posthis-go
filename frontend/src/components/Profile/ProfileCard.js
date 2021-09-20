import React from 'react';

import {
  Button
} from '@material-ui/core';

import { createFollow, deleteFollow }               from '_api';
import { authenticationService, followService}      from '_services';
import { FollowViewModel }                          from '_model';

import profilePicPlaceholder from 'assets/avatar-placeholder.svg';

import styles from '_styles/ProfileCard.module.css';

export const ProfileCard = ( props ) => {
  const { user, setUser } = props;

  const onClickFollowButton = async ()=>{
    const model = new FollowViewModel({
      followedID: user.id,
      followerID: authenticationService.currentUserValue.id });
    
    if( user.isFollowed ){
      const {data:responseData, err} = await deleteFollow( user.id );

      if ( err !== null ) return;

      const {data} = responseData;
      setUser(data);
    }else{
      const {data:responseData, err} = await createFollow( user.id );

      if ( err !== null ) return;

      const {data} = responseData;
      setUser(data);
    }
  };

  return (
    <div className={styles.cardContainer}>
      <div className={styles.card}>
          
        <div className={styles.contImg}>
          <img className={styles.profilePicture} src={ user.profilePicPath || profilePicPlaceholder }/>
          <h2 className={styles.username}>
            {user.username}
          </h2>
          <h4 className={styles.tag}>
            {`@${user.tag}`}
          </h4>
          <div className={styles.follows}>
            <span 
              className = { styles.followLink }
              onClick   = { ()=>followService.getFollowerUsers(user.id) }
            >
              {user.followerCount} Followers
            </span>
            <span 
              className = { styles.followLink }
              onClick   = { ()=>followService.getFollowedUsers(user.id) }
            >
              {user.followingCount} Following
            </span>
          </div>
        </div>

        {
          (authenticationService.currentUserValue.id !== user.id) 
          &&
          <div className={styles.followBtnContainer}>
            <Button
              fullWidth
              variant='contained'
              color='secondary'
              className={styles.followBtn}
              onClick={onClickFollowButton}
            >
              { user.isFollowed ? 'Unfollow' : 'Follow'}
            </Button>
          </div>
        }

        
      </div>
    </div>
  );
}

export default ProfileCard;
