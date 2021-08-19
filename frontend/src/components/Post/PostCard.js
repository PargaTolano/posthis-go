import React, { useState, useRef }  from 'react';
import { Link }                     from 'react-router-dom';

import {
  CardActions,
  CardContent,
  Typography,
  IconButton,
} from '@material-ui/core';

import {
  Favorite as FavoriteIcon,
  QuestionAnswer as QuestionAnswerIcon,
  ReplyAll       as ReplyAllIcon,
  Save           as SaveIcon,
  Image          as ImageIcon,
  Edit           as EditIcon,
  Delete         as DeleteIcon,
  Cancel         as CancelIcon,
} from '@material-ui/icons';

import { MediaGrid }                from 'components/Media';

import { handleResponse, prettyDate }           from '_helpers';
import { authenticationService }    from '_services';
import { routes, fileToBase64 }     from '_utils';

import { 
  updatePost,
  deletePost,
  createLike,
  deleteLike,
  createRepost,
  deleteRepost
} from '_api';

import { UPostModel, RepostViewModel } from '_model';

import Placeholder from 'assets/avatar-placeholder.svg'

import styles from '_styles/PostCard.module.css';

export const PostCard = ( props ) => {

  const { post:postProp, history } = props;

  const [ post, setPost ] = useState( postProp );

  const [ state, setState ] = useState({
    editMode:         false,
    content:          post.content,
    originalContent:  post.content,
    medias:           post.medias,
    originalMedias:   post.medias,
    deleted:          [],
    newMedias:        []
  });

  const inputFileRef = useRef();

  const temp = { state, setState };

  const onClickSave = async ()=>{
    const {data:responseData, err} =
      await updatePost( 
        post.postID,  
        new UPostModel({
          content: state.content, 
          deleted: state.deleted, 
          files: state.newMedias})
      );

    if (err !== null) return;

    const {data} = responseData

    setState(x=>{
      let copy              = {...x};
      copy.editMode         = false;
      copy.originalContent  = copy.content;
      copy.medias           = data.medias;
      copy.originalMedias   = data.medias;
      return copy;
    });
  };
  
  const onChangeContent = e=>{
      setState( x =>{
        let copy = {...x};
        copy.content = e.target.value;
        return copy;
      });
  };

  const onChangeImages = async e=>{
    let {files} = e.target;

    if( state.medias?.length + files.length > 4 )
      return;

    const mediaInfo    = [];
    const newMediaFiles = [];
    for ( let i = 0; i < files.length; i++ ){
      let file = files[i];
      let preview = await fileToBase64( file );

      const mediaViewModel = {
        mediaID:  null,
        mime:     file.type,
        path:     preview,
        isVideo:  file.type.includes('video')
      };

      mediaInfo     .push( mediaViewModel );
      newMediaFiles .push( file );
    }

    setState( x=>{
      let copy = {...x};

      if (copy.medias)
        copy.medias = [...copy.medias, ...mediaInfo];
      else 
        copy.medias = [...mediaInfo];

      if (copy.medias)
        copy.newMedias = [...copy.newMedias, ...newMediaFiles];
      else 
        copy.medias = [...newMediaFiles];

      return copy;
    });
  };

  const onClickFileOpen = ()=>inputFileRef.current?.click();

  const onClickLike = async ()=>{

    if( post.isLiked ){

      const {data:responseData, err}  = await deleteLike(post.postID);

      if ( err != null ) return;

      const { data } = responseData;

      if (data === null ) return;

      setPost( data );
      setState( x=>{
        let copy = {...x};
        copy.content         = data.content;
        copy.originalContent = data.content;
        copy.medias          = data.medias;
        copy.originalMedias  = data.medias;
        return copy;
      });
    }
    else{

      const {data:responseData, err}  = await createLike(authenticationService.currentUserValue.id, post.postID);

      if ( err != null ) return;

      const { data } = responseData;

      if (data === null ) return;

      setPost( data );
      setState( x=>{
        let copy = {...x};
        copy.content         = data.content;
        copy.originalContent = data.content;
        copy.medias          = data.medias;
        copy.originalMedias  = data.medias;
        return copy;
      });
    }
  };

  const onClickRepost = async ()=>{

    if( post.isReposted ) {
      const {data:responseData, err} = await deleteRepost(post.postID);

      if ( err !== null) return;

      const { data } = responseData;

      setPost( data );
      setState( x=>{
        let copy = {...x};
        copy.content         = data.content;
        copy.originalContent = data.content;
        copy.medias          = data.medias;
        copy.originalMedias  = data.medias;
        return copy;
      });
    }
    else{

      const {data:responseData, err} = await createRepost(authenticationService.currentUserValue.id, post.postID);

      if ( err !== null) return;

      const { data } = responseData;

      setPost( data );
      setState( x=>{
        let copy = {...x};
        copy.content         = data.content;
        copy.originalContent = data.content;
        copy.medias          = data.medias;
        copy.originalMedias  = data.medias;
        return copy;
      });
    }
  };

  const onToggleEditMode = ()=>{
    setState( x =>{
      let copy = {...x};
      copy.editMode = !copy.editMode;

      if( copy.editMode === false ){

        if( copy.originalMedias)
          copy.medias   = [...copy.originalMedias];
        else
          copy.medias   = null;
        copy.content    = copy.originalContent;
        copy.deleted    = [];
        copy.newMedias  = [];
      }
      return copy;
    })
  };

  const onClickDelete = async ()=>{
    if(window.confirm('Seguro que quiere borrar el post')){
      const {err} = await deletePost(post.postID);

      if ( err !== null ) return;

      history.replace(routes.feed);
    }
  };

  return (
    <div className={styles.root}>
      {
        (!!post.isRepost) &&
        <Typography variant='body2' className={styles.repostText}>
          <Link to={routes.getProfile(post.publisherID)} className={styles.repostUserLink}>  
            <ReplyAllIcon className={styles.repostedIcon}/> {post.reposterUserName} reposted this!
          </Link> 
        </Typography>
      }
      <CardContent>
          <div className={styles.displayTitle}>
            <Link to={routes.getProfile(post.publisherID)} className={styles.avatarContainer}>
              <img src={post.publisherProfilePic || Placeholder} className={styles.avatar}/>
            </Link>
            <Link to={routes.getProfile(post.publisherID)} className={styles.titleContainer}>
              <Typography variant='h6' component='h2' className={styles.title}>
                <strong className={styles.publisher}>{post.publisherUserName} {'@'+post.publisherTag}</strong>
                <p className={styles.date}> { prettyDate(post.date) } </p>
              </Typography>
            </Link>
            <div  className={styles.displaybtn}>
              {
                (authenticationService.currentUserValue.id === post.publisherID)
                && 
                <>
                  <IconButton 
                    variant='contained' 
                    color='secondary' 
                    onClick={onToggleEditMode}
                  >
                    { state.editMode ?  <CancelIcon className={styles.cancelIcon}/> : <EditIcon className={styles.saveIcon}/>  }
                  </IconButton>
                  { 
                    !state.editMode && <IconButton 
                      variant='contained' 
                      color='secondary' 
                      onClick={onClickDelete}
                      className={styles.deleteBtn}
                    >
                      <DeleteIcon className={styles.cancelIcon}/> 
                    </IconButton>
                    }
                  {
                    (authenticationService.currentUserValue.id === post.publisherID && state.editMode)
                    && 
                    <div>
                      <IconButton onClick={onClickSave}>
                        <SaveIcon className={styles.saveIcon}/>
                      </IconButton>
                    </div>
                  }
                </>
              }
            </div>
          </div>


          {
            state.editMode ? 
            (
              <textarea className={styles.contentEdit} value={state.content} onChange={onChangeContent}></textarea>
            )
            :
            (
              <Typography variant='body2' component='p' className={ ( ( state.medias === null && ( state.originalContent.split('\n').length < 5 ) ) ) ? styles.contentNoMedia : styles.content} >
                {state.originalContent}
              </Typography>
            )
          }
          
          {
            state.editMode &&
            <>
              <input accept='image/*' className={styles.input} type='file' multiple ref={inputFileRef} onChange={onChangeImages}/>
              <IconButton onClick={onClickFileOpen}>
                <ImageIcon className={styles.mediaIcon}/>
              </IconButton>
            </>
          }
          <div className={styles.contMedia}>
            <MediaGrid media={ state.editMode ? state.medias : state.originalMedias } {...temp}/>
          </div>

          
        </CardContent>

      <CardActions disableSpacing className={styles.cardBtn}>

        {
          (!state.editMode)
          && 
          <div>
            <IconButton onClick={onClickLike}>
              <FavoriteIcon className={ post.isLiked ? styles.likeIcon : styles.grayIcon }/>
            </IconButton>
            {post.likeCount}
          </div>
        }
        
        {
          (!state.editMode)
          && 
          <div>
            <Link to={routes.getPost(post.postID)}>
              <IconButton>
                <QuestionAnswerIcon className={styles.commentIcon}/>
              </IconButton>
            </Link>
            {post.replyCount}
          </div>
        }
        
        {
          (!state.editMode)
          && 
          <div>
            <IconButton onClick={onClickRepost}>
              <ReplyAllIcon className={ post.isReposted ? styles.repostIcon : styles.grayIcon }/>
            </IconButton>
            {post.repostCount}
          </div>
        }

      </CardActions>
    </div>
  );
}

export default PostCard;