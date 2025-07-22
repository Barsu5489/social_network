
'use client';

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Logo } from "@/components/logo";
import { useRouter } from 'next/navigation';
import { useForm, UseFormReturn } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { cn } from "@/lib/utils";
import { Eye, EyeOff } from 'lucide-react';

import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Switch } from "@/components/ui/switch";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/hooks/use-toast";
import { API_BASE_URL } from "@/lib/config";
import { useUser } from "@/contexts/user-context";
import { FileUpload } from "@/components/file-upload";
import { getRandomAvatar } from "@/lib/avatars";

const signUpSchema = z.object({
  firstName: z.string().min(1, { message: "First name is required." }),
  lastName: z.string().min(1, { message: "Last name is required." }),
  email: z.string().email({ message: "Please enter a valid email." }),
  password: z.string().min(8, { message: "Password must be at least 8 characters." }),
  dateOfBirth: z.string()
    .min(1, { message: "Your date of birth is required." })
    .refine(
      (val) => {
        if (!/^\d{4}-\d{2}-\d{2}$/.test(val)) return false;
        const date = new Date(val);
        if (isNaN(date.getTime())) return false;
        if (date > new Date() || date < new Date("1900-01-01")) return false;
        return true;
      },
      {
        message: "Please enter a valid date in YYYY-MM-DD format (e.g., 1990-01-31).",
      }
    ),
  isPrivate: z.boolean().default(false),
  nickname: z.string().optional(),
  aboutMe: z.string().max(200, { message: "About me must be 200 characters or less." }).optional(),
  avatarFile: z.any().optional().refine(
    (file) => !file || (typeof File !== 'undefined' && file instanceof File),
    { message: "Invalid file type" }
  ),
});

const signInSchema = z.object({
    email: z.string().email({ message: "Please enter a valid email." }),
    password: z.string().min(1, { message: "Password is required" }),
});

type SignUpFormValues = z.infer<typeof signUpSchema>;
type SignInFormValues = z.infer<typeof signInSchema>;

// SignInForm Component
const SignInFormComponent = ({
  form,
  onSubmit,
  showPassword,
  setShowPassword,
  onSwitchMode,
}: {
  form: UseFormReturn<SignInFormValues>;
  onSubmit: (values: SignInFormValues) => void;
  showPassword: boolean;
  setShowPassword: React.Dispatch<React.SetStateAction<boolean>>;
  onSwitchMode: () => void;
}) => (
  <>
    <div className="grid gap-2 text-center">
      <h1 className="text-3xl font-bold">Sign In</h1>
      <p className="text-balance text-muted-foreground">Enter your email below to login to your account</p>
    </div>
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="grid gap-4">
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
              <FormItem>
                  <FormLabel>Email</FormLabel>
                  <FormControl>
                      <Input type="email" placeholder="m@example.com" {...field} />
                  </FormControl>
                  <FormMessage />
              </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
              <FormItem>
                  <FormLabel>Password</FormLabel>
                  <div className="relative">
                    <FormControl>
                        <Input type={showPassword ? "text" : "password"} {...field} />
                    </FormControl>
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 text-muted-foreground"
                      onClick={() => setShowPassword((prev) => !prev)}
                    >
                      {showPassword ? <EyeOff className="h-4 w-4"/> : <Eye className="h-4 w-4"/>}
                      <span className="sr-only">Toggle password visibility</span>
                    </Button>
                  </div>
                  <FormMessage />
              </FormItem>
          )}
        />
        <Button type="submit" className="w-full" disabled={form.formState.isSubmitting}>
          {form.formState.isSubmitting ? "Signing In..." : "Sign In"}
        </Button>
      </form>
    </Form>
    <div className="mt-4 text-center text-sm">
      Don&apos;t have an account?{' '}
      <button onClick={onSwitchMode} className="underline font-semibold text-primary">
        Sign up
      </button>
    </div>
  </>
);

// SignUpForm Component (Multi-step)
const SignUpFormComponent = ({
  form,
  onSubmit,
  showPassword,
  setShowPassword,
  onSwitchMode,
  step,
  setStep,
}: {
  form: UseFormReturn<SignUpFormValues>;
  onSubmit: (values: SignUpFormValues) => void;
  showPassword: boolean;
  setShowPassword: React.Dispatch<React.SetStateAction<boolean>>;
  onSwitchMode: () => void;
  step: number;
  setStep: React.Dispatch<React.SetStateAction<number>>;
}) => {
  const handleNext = async () => {
    const fields: (keyof SignUpFormValues)[] = ['firstName', 'lastName', 'email', 'password', 'dateOfBirth'];
    const isValid = await form.trigger(fields, { shouldFocus: true });
    if (isValid) {
      setStep(2);
    }
  };

  return (
    <>
      <div className="grid gap-2 text-center">
        <h1 className="text-3xl font-bold">Create an account</h1>
        <p className="text-balance text-muted-foreground">
          {step === 1 ? "Enter your essential details to get started." : "Add some optional flair to your profile."}
        </p>
      </div>
      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="grid gap-4">
          {step === 1 && (
            <>
              <div className="grid grid-cols-2 gap-4">
                <FormField control={form.control} name="firstName" render={({ field }) => (<FormItem><FormLabel>First Name</FormLabel><FormControl><Input placeholder="Sofia" {...field} /></FormControl><FormMessage /></FormItem>)} />
                <FormField control={form.control} name="lastName" render={({ field }) => (<FormItem><FormLabel>Last Name</FormLabel><FormControl><Input placeholder="Robinson" {...field} /></FormControl><FormMessage /></FormItem>)} />
              </div>
              <FormField control={form.control} name="email" render={({ field }) => (<FormItem><FormLabel>Email</FormLabel><FormControl><Input type="email" placeholder="m@example.com" {...field} /></FormControl><FormMessage /></FormItem>)} />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                    <FormItem>
                        <FormLabel>Password</FormLabel>
                        <div className="relative">
                          <FormControl>
                              <Input type={showPassword ? "text" : "password"} {...field} />
                          </FormControl>
                          <Button type="button" variant="ghost" size="icon" className="absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7 text-muted-foreground" onClick={() => setShowPassword((prev) => !prev)}>
                            {showPassword ? <EyeOff className="h-4 w-4"/> : <Eye className="h-4 w-4"/>}
                            <span className="sr-only">Toggle password visibility</span>
                          </Button>
                        </div>
                        <FormMessage />
                    </FormItem>
                )}
              />
              <FormField control={form.control} name="dateOfBirth" render={({ field }) => (<FormItem><FormLabel>Date of birth</FormLabel><FormControl><Input placeholder="YYYY-MM-DD" {...field} /></FormControl><FormMessage /></FormItem>)} />
              <Button type="button" className="w-full" onClick={handleNext}>Next</Button>
            </>
          )}

          {step === 2 && (
            <>
              <FormField control={form.control} name="nickname" render={({ field }) => (<FormItem><FormLabel>Nickname (Optional)</FormLabel><FormControl><Input placeholder="Your cool nickname" {...field} /></FormControl><FormMessage /></FormItem>)} />
              <FormField control={form.control} name="aboutMe" render={({ field }) => (<FormItem><FormLabel>About Me (Optional)</FormLabel><FormControl><Textarea placeholder="Tell us a little about yourself" {...field} /></FormControl><FormMessage /></FormItem>)} />
              <FormField
                  control={form.control}
                  name="avatarFile"
                  render={({ field }) => (
                      <FormItem>
                          <FormLabel>Avatar (Optional)</FormLabel>
                          <FormControl>
                              <FileUpload
                                  value={field.value}
                                  onChange={field.onChange}
                                  disabled={form.formState.isSubmitting}
                              />
                          </FormControl>
                          <FormMessage />
                      </FormItem>
                  )}
              />
              <FormField control={form.control} name="isPrivate" render={({ field }) => (<FormItem className="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm"><div className="space-y-0.5"><FormLabel>Private Account</FormLabel><FormDescription>Only followers will see your posts.</FormDescription></div><FormControl><Switch checked={field.value} onCheckedChange={field.onChange} /></FormControl></FormItem>)} />
              
              <div className="grid grid-cols-2 gap-4">
                <Button type="button" variant="outline" onClick={() => setStep(1)}>Back</Button>
                <Button type="submit" disabled={form.formState.isSubmitting}>
                  {form.formState.isSubmitting ? "Creating Account..." : "Create Account"}
                </Button>
              </div>
            </>
          )}
        </form>
      </Form>
      <div className="mt-4 text-center text-sm">
        Already have an account?{' '}
        <button onClick={onSwitchMode} className="underline font-semibold text-primary">
          Sign in
        </button>
      </div>
    </>
  );
};


export default function AuthPage() {
  const [authMode, setAuthMode] = useState<'signin' | 'signup'>('signin');
  const [signupStep, setSignupStep] = useState(1);
  const [showSignInPassword, setShowSignInPassword] = useState(false);
  const [showSignUpPassword, setShowSignUpPassword] = useState(false);
  const router = useRouter();
  const { toast } = useToast();
  const { setUser } = useUser();

  const signUpForm = useForm<SignUpFormValues>({
    resolver: zodResolver(signUpSchema),
    defaultValues: {
      firstName: "",
      lastName: "",
      email: "",
      password: "",
      dateOfBirth: "",
      isPrivate: false,
      nickname: "",
      aboutMe: "",
    },
  });

  const signInForm = useForm<SignInFormValues>({
    resolver: zodResolver(signInSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });
  
  const handleSwitchMode = () => {
    setAuthMode(authMode === 'signin' ? 'signup' : 'signin');
    setSignupStep(1); // Reset to first step
  };


  async function handleSignIn(values: SignInFormValues) {
    try {
        const response = await fetch(`${API_BASE_URL}/api/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(values),
            credentials: 'include',
        });

        if (response.ok) {
            const data = await response.json();
            setUser(data.user); // Change from data.data to data.user
            router.push('/home');
        } else {
            let errorMessage = 'Please check your credentials and try again.';
            try {
                const errorData = await response.json();
                errorMessage = errorData.error || errorData.message || errorMessage;
            } catch (jsonError) {
                try {
                    const errorText = await response.text();
                    if (errorText) {
                        errorMessage = errorText;
                    }
                } catch (textError) {
                    errorMessage = `An error occurred: ${response.statusText}`;
                }
            }
            toast({
                variant: "destructive",
                title: "Sign In Failed",
                description: errorMessage,
            });
        }
    } catch (error) {
        toast({
            variant: "destructive",
            title: "Network Error",
            description: "Could not connect to the server. Please try again later.",
        });
    }
  }

  async function handleSignUp(values: SignUpFormValues) {
    const { firstName, lastName, dateOfBirth, isPrivate, aboutMe, avatarFile, ...payload } = values;
    
    // NOTE: This implementation doesn't upload the file. A real-world scenario would
    // involve uploading the `avatarFile` to a storage service (like Cloud Storage)
    // and then passing the returned URL to the backend. For now, we send an empty URL.
    const requestBody = {
      ...payload,
      first_name: firstName,
      last_name: lastName,
      date_of_birth: dateOfBirth,
      is_private: isPrivate,
      about_me: aboutMe,
      avatar_url: avatarFile ? "" : getRandomAvatar(),
    };
    
    try {
        const response = await fetch(`${API_BASE_URL}/api/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody),
            credentials: 'include',
        });

        if (response.ok) {
            toast({
                title: "Registration Successful",
                description: "Welcome! Please sign in to continue.",
            });
            await handleSignIn({ email: values.email, password: values.password });
        } else {
            let errorMessage = 'Please check your details and try again.';
            try {
                const errorData = await response.json();
                errorMessage = errorData.error || errorData.message || errorMessage;
            } catch (jsonError) {
                try {
                    const errorText = await response.text();
                    if (errorText) {
                        errorMessage = errorText;
                    }
                } catch (textError) {
                    errorMessage = `An error occurred: ${response.statusText}`;
                }
            }
            toast({
                variant: "destructive",
                title: "Sign Up Failed",
                description: errorMessage,
            });
        }
    } catch (error) {
        toast({
            variant: "destructive",
            title: "Network Error",
            description: "Could not connect to the server. Please try again later.",
        });
    }
  }

  return (
    <div className="w-full min-h-screen lg:grid lg:grid-cols-2">
      <div className="flex items-center justify-center py-12 px-4 sm:px-0">
        <div className="mx-auto w-full max-w-[400px] grid gap-6 bg-background/80 backdrop-blur-sm border rounded-xl p-8 shadow-2xl">
          {authMode === 'signin' ? (
            <SignInFormComponent
              form={signInForm}
              onSubmit={handleSignIn}
              showPassword={showSignInPassword}
              setShowPassword={setShowSignInPassword}
              onSwitchMode={handleSwitchMode}
            />
          ) : (
            <SignUpFormComponent
              form={signUpForm}
              onSubmit={handleSignUp}
              showPassword={showSignUpPassword}
              setShowPassword={setShowSignUpPassword}
              onSwitchMode={handleSwitchMode}
              step={signupStep}
              setStep={setSignupStep}
            />
          )}
        </div>
      </div>
      <div className="hidden lg:flex items-center justify-center p-8 flex-col bg-muted/40">
        <Logo />
        <div className="text-center mt-6">
            <h2 className="text-4xl font-bold font-headline">Connect with your world.</h2>
            <p className="mt-2 text-lg text-muted-foreground">Join a community of creators, friends, and innovators.</p>
        </div>
      </div>
    </div>
  );
}
