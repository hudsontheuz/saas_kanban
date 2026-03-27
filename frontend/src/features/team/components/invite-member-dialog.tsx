import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

export function InviteMemberDialog({ onInvite, disabled }: { onInvite: (name: string, email: string) => Promise<void>; disabled?: boolean }) {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');

  const handleInvite = async () => {
    await onInvite(name, email);
    setName('');
    setEmail('');
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button disabled={disabled}>Convidar membro</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Novo membro</DialogTitle>
          <DialogDescription>Limite atual da equipe: 5 pessoas.</DialogDescription>
        </DialogHeader>
        <div className="space-y-4 pt-4">
          <div className="space-y-2">
            <Label htmlFor="invite-name">Nome</Label>
            <Input id="invite-name" value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="invite-email">E-mail</Label>
            <Input id="invite-email" value={email} onChange={(e) => setEmail(e.target.value)} />
          </div>
          <Button className="w-full" onClick={handleInvite}>Adicionar</Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
