
'use client';

import { UploadCloud, X } from 'lucide-react';
import React, { useCallback, useState, useEffect } from 'react';
import { useDropzone, FileRejection, DropzoneOptions } from 'react-dropzone';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import Image from 'next/image';
import { useToast } from '@/hooks/use-toast';

interface FileUploadProps {
  onChange: (file: File | undefined) => void;
  value?: File | null;
  disabled?: boolean;
}

export function FileUpload({ onChange, value, disabled }: FileUploadProps) {
  const { toast } = useToast();
  const [preview, setPreview] = useState<string | null>(null);

  useEffect(() => {
    if (value) {
      const objectUrl = URL.createObjectURL(value);
      setPreview(objectUrl);
      return () => URL.revokeObjectURL(objectUrl);
    }
    setPreview(null);
  }, [value]);

  const onDrop = useCallback(
    (acceptedFiles: File[], fileRejections: FileRejection[]) => {
      if (fileRejections.length > 0) {
        fileRejections.forEach(({ errors }) => {
          errors.forEach((error) => {
            toast({
              title: "Upload Error",
              description: error.message,
              variant: "destructive",
            });
          });
        });
        return;
      }
      
      const file = acceptedFiles[0];
      if (file) {
        onChange(file);
      }
    },
    [onChange, toast]
  );

  const dropzoneOptions: DropzoneOptions = {
    onDrop,
    accept: { 'image/png': ['.png'], 'image/jpeg': ['.jpg', '.jpeg'], 'image/gif': ['.gif'] },
    maxSize: 1024 * 1024, // 1MB
    multiple: false,
    disabled,
  };

  const { getRootProps, getInputProps, isDragActive } = useDropzone(dropzoneOptions);

  const removeFile = (e: React.MouseEvent<HTMLButtonElement>) => {
    e.stopPropagation();
    onChange(undefined);
  };

  return (
    <div
      {...getRootProps()}
      className={cn(
        'relative flex flex-col items-center justify-center w-full h-48 border-2 border-dashed rounded-lg cursor-pointer bg-secondary/50 hover:bg-secondary/80 transition-colors',
        isDragActive && 'border-primary bg-primary/10',
        disabled && 'cursor-not-allowed opacity-50',
        preview && 'border-solid p-0 overflow-hidden'
      )}
    >
      <input {...getInputProps()} />
      {preview ? (
        <>
          <Image src={preview} alt="Avatar preview" fill style={{ objectFit: 'contain' }} className="rounded-lg p-2" />
          {!disabled && (
            <Button
              type="button"
              variant="destructive"
              size="icon"
              className="absolute top-2 right-2 h-6 w-6 z-10"
              onClick={removeFile}
            >
              <X className="h-4 w-4" />
            </Button>
          )}
        </>
      ) : (
        <div className="flex flex-col items-center justify-center text-center text-muted-foreground p-4">
          <UploadCloud className="h-10 w-10 mb-2" />
          <p className="font-semibold">Drag & drop image here, or click to select</p>
          <p className="text-xs mt-1">PNG, JPG, GIF up to 1MB</p>
        </div>
      )}
    </div>
  );
}
