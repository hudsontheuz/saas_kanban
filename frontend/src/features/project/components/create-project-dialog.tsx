import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

export function CreateProjectDialog({ onCreate, disabled }: { onCreate: (name: string) => Promise<void>; disabled?: boolean }) {
  const [name, setName] = useState('');

  const handleCreate = async () => {
    await onCreate(name);
    setName('');
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button disabled={disabled}>Novo projeto</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Criar projeto</DialogTitle>
          <DialogDescription>Dê um nome ao projeto que ficará ativo para a equipe.</DialogDescription>
        </DialogHeader>
        <div className="space-y-4 pt-4">
          <div className="space-y-2">
            <Label htmlFor="project-name">Nome do projeto</Label>
            <Input id="project-name" value={name} onChange={(e) => setName(e.target.value)} placeholder="Ex.: Lançamento do MVP" />
          </div>
          <Button className="w-full" onClick={handleCreate} disabled={!name.trim()}>Salvar</Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
