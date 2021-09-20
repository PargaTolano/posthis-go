import React from 'react';

import PostCard       from 'components/Post/PostCard';

import styles from '_styles/PostContainer.module.css'

export const PostContainer = (props) => {

  const {  posts, history } = props;

  return (
    <div className={styles.cardHolder}>
      {
        posts 
        &&
        posts.map(x=><PostCard key={x.postID} post={x} history={history}/>)
      }
        
    </div>
  );
}
export default PostContainer;
