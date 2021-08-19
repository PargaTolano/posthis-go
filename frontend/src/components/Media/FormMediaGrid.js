import React from 'react';

import { makeStyles } from '@material-ui/core/styles';

import {
    IconButton
} from '@material-ui/core';

import {
    Cancel as CancelIcon
} from '@material-ui/icons';

const useStyles = makeStyles((theme)=>({
    previewGrid: {
        display:        'flex',
        width:          '100%',
        justifyContent: 'center',
        flexWrap:       'wrap',
        borderRadius:   '10px',
        boxShadow:      'black 0px 0px 2px',
        overflow:       'hidden',
    },
    previewContainer:{
        display:    'inline-block',
        flexGrow:   '1',
        width:      '50%',
        height:     '180px',
    },
    previewImage:{
        postion:    'absolute',
        top:        '0',
        left:       '0',
        display:    'inline-block',
        width:      '100%',
        height:     '100%',
        objectFit:  'cover'
    },
    closePreviewIcon:{
        position: 'absolute',
        zIndex: 1,
    }
}));

const GridImage = ( props ) => {

  const { classes, image, index, images, setImages } = props;
  const { preview, file } = image;

  const deleteImage = ()=>{
    const i = images.indexOf( image );
    setImages( x => x.filter( (elem,index)=> index != i ) );
  };

  return (
    <div className={classes.previewContainer}>
      <IconButton className={classes.closePreviewIcon} color='secondary' aria-label='upload picture' component='span' onClick={deleteImage}>
        <CancelIcon/>
      </IconButton>
      <img src={preview} className={classes.previewImage}/>
    </div>
  );
};

export const FormMediaGrid = (props) =>{

    const classes = useStyles();
    const { images, setImages } = props;

    return (
        <div className={classes.previewGrid}>
            {
              images.map( (image, i)=>( <GridImage
                                          key={i}
                                          index={i}
                                          image={image}
                                          images={images}
                                          setImages={setImages}
                                          classes={classes}
                                        />))
            }
      </div>
  );
};

export default FormMediaGrid;