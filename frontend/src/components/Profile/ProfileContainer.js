import React, { useState, useEffect } from 'react';

import { makeStyles }                 from '@material-ui/core/styles';

import { ProfileCard }                from 'components/Profile';
import { PostContainer }              from 'components/Feed';
import { CreatePostForm }             from 'components/Post';
import { PaginationElement }          from 'components/Common';

import { authenticationService }      from '_services';
import { getUserFeed }                from '_api';
import { loadingState }               from '_hooks';

import { DialogEditInfo, EditInfo }   from 'components/EditProfile';

import coverPlaceholder from 'assets/background-placeholder.jpg';

import styles from '_styles/ProfileContainer.module.css';

export const ProfileContainer = ( props ) => {
  const { user, setUser, ...rest } = props;

  const [posts, setPosts] = useState([]);

  const loadFeed = async ()=>{

    loadingState.set(true);

    const {data:responseData, err} = await getUserFeed(user.id, 0, 100);

    loadingState.set(false);

    if (err !== null) 
      return;
    
    const { data } = responseData;
    setPosts( data );
  };

  useEffect(() => {
    loadFeed();
  },[user]);

  return (
      <div className={styles.root}>


        <div className={styles.coverContainer}>
          <img className={styles.coverPic} src={ user.coverPicPath || coverPlaceholder } />
          {
            (authenticationService.currentUserValue.id === user.id) 
            &&
            <DialogEditInfo color={'primary'}>
              <EditInfo user={user} setUser={setUser}/>
            </DialogEditInfo>
          }
          <ProfileCard user={user} setUser={setUser} {...rest}/>
        </div>


        {
          ( authenticationService.currentUserValue.id === user.id ) && <CreatePostForm afterUpdate={loadFeed} {...rest}/>
        }

        <PostContainer posts={posts} {...rest}/>

        <PaginationElement
          name            = 'posts'
          hasFetched      = {true}
          total           = {0}
          last            = {0}
          limit           = {1}
          onIntersection  = {()=>{}}
        />
      </div>
  );
}
export default ProfileContainer;