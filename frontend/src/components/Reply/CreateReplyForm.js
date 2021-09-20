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

import { FormMediaGrid }                            from 'components/Media';

import { handleResponse }                           from '_helpers';
import { authenticationService, replyService }      from '_services';
import { fileToBase64,validateCreateAndUpdatePost } from '_utils';
import { createReply }    from '_api';
import {CReplyModel}      from '_model';

const useStyles = makeStyles((theme) => ({
  form: {
    background:     'white',
    width:          '100%',
    maxWidth:       '700px',
    height:         '40%',
    padding:        theme.spacing( 3 ),
    borderRadius:   theme.spacing( 0, 0, 1, 1),  
    marginBottom:   theme.spacing( 2 ),
    boxShadow:      'black 1px 1px 8px',
    alignItems:     'center',
    justifySelf:    'center',
    alignSelf:      'center',
    marginLeft:     'auto',
    marginRight:    'auto',
    zIndex:         1,
  },
  cardBtn: {
    display: 'flex',
    flexWrap: 'wrap',
    alignItems: 'center',
    justifyContent: 'space-between',
    color: 'white',
    alignSelf:'center',
  },
  submit: {
    order: 0,
    width: '20%',
    [theme.breakpoints.down('sm')]:{
      order: 1,
      width: '100%',
      flexGrow: '1',
      flexShrink: '0',
      flexBasis: 'auto',
    }
  },
  imageIcon: {
    order: 1,
    [theme.breakpoints.down('sm')]:{
      order: 0,
      flexGrow: '1',
      flexShrink: '0',
      flexBasis: 'auto',
    }
  },
  titleForm: {
    justifyContent: 'space-around',
    color: theme.palette.primary.dark ,
    alignSelf:'center',
    color:'#ea5970',
  },
  input: {
    display: 'none',
  },
  
}));

export const CreateReplyForm = (props) => {

  const {postId} = props;
  
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
  
  const classes = useStyles();
  
  const onSubmit = async ( e )=>{
    e.preventDefault();

    const model = new CReplyModel({
        userID: authenticationService.currentUserValue.id,
        postID: postId, 
        content,
        files: images.map(x=>x.file)});

    const { err } = await createReply( model );

    if(err != null) 
      return;

    replyService.getPostReplies(postId);

    setImages([]);
    setContent('');
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
    <form 
      className={classes.form} 
      noValidate 
      onSubmit={onSubmit}
    >
        <div 
          component='h4' 
          variant='h2' 
          className={classes.titleForm}
        >
          <strong>Reply</strong>
        </div>
        <TextField
            variant      = 'outlined'
            margin       = 'normal'
            multiline
            rows         = {3}
            rowsMax      = {3}
            fullWidth
            label        = "What's do you think?"
            name         = 'postContent'
            autoComplete = 'postContent'
            autoFocus
            value        = {content}
            className    = {classes.postContent}
            onChange     = {onChangeContent}
        />

      <FormMediaGrid images={images} setImages={setImages}/>
      <div className = {classes.cardBtn}>
      
        <Button
            type     = 'submit'
            fullWidth
            variant  = 'contained'
            color    = 'primary'
            className= {classes.submit}
            disabled = {!validation.validated}
        >
          Reply!
        </Button>
          
        <input 
          accept='image/*' 
          className={classes.input} 
          type='file' 
          multiple 
          ref={inputFileRef} 
          onChange={onChangeImage}
        />
        <label
          htmlFor='icon-button-file' 
          className={classes.imageIcon}
        >
          <IconButton
            color='primary' 
            aria-label='upload picture' 
            component='span' 
            onClick={mediaBtnOnClick}
          >
            <ImageIcon className={classes.imageIcon}/>
          </IconButton>
        </label>
        
        </div>
    </form>
  );
}
export default CreateReplyForm;