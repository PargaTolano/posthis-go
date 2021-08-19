import React,{useState, useEffect, useRef} from 'react';

import { makeStyles } from '@material-ui/core/styles';

import {
  Image as ImageIcon
} from '@material-ui/icons';

import {
  TextField,
  Button,
  IconButton,
} from '@material-ui/core';

import { FormMediaGrid }          from 'components/Media';

import { authenticationService }  from '_services';
import { 
  fileToBase64,
   validateCreateAndUpdatePost 
} from '_utils';
import { createPost }             from '_api';

import CPostModel     from '_model/CPostModel';

import styles from '_styles/PostForm.module.css';

export const CreatePostForm = (props) => {

  const {afterUpdate} = props;
  
  const [images, setImages]   = useState( [] );
  const [content, setContent] = useState( '' );

  const [validation, setValidation] = useState({
    content:    false,
    mediaCount: false,
    validated:  false
  });

  useEffect(() => {
    setValidation( validateCreateAndUpdatePost( {content, mediaCount: images.length}) );
  }, [images, content])

  const inputFileRef = useRef(null);
  
  const onSubmit = async ( e )=>{
    e.preventDefault();
    const { err } = await createPost( new CPostModel({
        userID: authenticationService.currentUserValue.id,
        content,
        files: images.map(x=>x.file)
      }));

    setImages([]);
    setContent('');
    if( afterUpdate && err === null)
      afterUpdate();
  };

  const onChangeImage = async ( e )=>{
    let { files } = e.target;

    if( images.length + files.length > 4 )
      return;

    const filePairs = [];
    for ( let i = 0; i < files.length; i++ ){
      let file = files[i];
      let preview = await fileToBase64( file );
      filePairs.push( {
        file,
        preview
      });
    }

    setImages( x=> [...x,...filePairs] );
  };

  const onChangeContent = e=>setContent(e.target.value);

  const mediaBtnOnClick = () =>inputFileRef.current?.click();

  return (

    <form className={styles.form} noValidate onSubmit={onSubmit}>
        <div component='h4' variant='h2' className={styles.titleForm}>
          <strong>New PosThis!</strong>
        </div>
        <TextField
          variant='outlined'
          margin='normal'
          multiline
          rows={3}
          rowsMax={3}
          fullWidth
          label="What's on your mind?"
          name='postContent'
          autoComplete='postContent'
          autoFocus
          value = {content}
          className={styles.postContent}
          onChange ={onChangeContent}
        />

      <FormMediaGrid images={images} setImages={setImages}/>
      <div className = {styles.cardBtn}>
      
        <Button
          type='submit'
          fullWidth
          variant='contained'
          color='primary'
          style={{width: 200}}
          className={styles.submit}
          disabled = {!validation.validated}
        >
          PosThis!
        </Button>
          
        <input
          accept='image/*' 
          className={styles.input} 
          type='file' 
          multiple 
          ref={inputFileRef} 
          onChange={onChangeImage}  
        />
        <label 
          htmlFor='icon-button-file' 
          className={styles.imageIcon}
        >
          <IconButton 
            color='primary' 
            aria-label='upload picture' 
            component='span' 
            onClick={mediaBtnOnClick}
          >
            <ImageIcon className={styles.imageIcon}/>
          </IconButton>
        </label>
        
        </div>
    </form>
       

  );
}
export default CreatePostForm;