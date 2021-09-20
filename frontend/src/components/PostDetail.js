import React                        from 'react';
import { Redirect }                 from 'react-router-dom';

import { NavBar }                   from 'components/Feed';

import {
  PostCard
} from 'components/Post';

import { 
  CreateReplyForm,
  ReplyContainer
} from 'components/Reply';


import { routes }                                 from '_utils';
import { useGetDetailedPost                       } from '_hooks';

import styles from '_styles/PostDetail.module.css';

export const PostDetail = ( props ) => {
  
  const { match, history, ...rest } = props;
  const { id }    = match.params;

  const [ready, post] = useGetDetailedPost( id );

  if( id == 'undefined' || id === undefined || id === null || id === '' || (ready && post === null) ){
    return (<Redirect to={routes.feed}/>);
  }

  return (
    <div className={styles.background}>
      <NavBar history={history}/>

      {
        (ready && post)
        &&
        <>
          <div component='h4' variant='h2' className={styles.titleBegin}>
            <strong>Detailed Post</strong>
          </div>
          <div className={styles.cardHolder}>
            <PostCard post={post} history={history}/>
            <CreateReplyForm postId={post?.postID}/>
            <ReplyContainer id={post.postID}/>
          </div>
        </>
      }
    </div>
  );
};

export default PostDetail;
