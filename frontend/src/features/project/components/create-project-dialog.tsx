import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

export function CreateProjectDialog({ onCreate, disabled }: { onCreate: (name: string, description: string) => Promise<void>; disabled?: boolean }) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  const handleCreate = async () => {
    await onCreate(name, description);
    setName('');
    setDescription('');
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button disabled={disabled}>Novo projeto</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Criar projeto</DialogTitle>
          <DialogDescription>Neste MVP a equipe mantém um projeto ativo por vez.</DialogDescription>
        </DialogHeader>
        <div className="space-y-4 pt-4">
          <div className="space-y-2">
            <Label htmlFor="project-name">Nome</Label>
            <Input id="project-name" value={name} onChange={(e) => setName(e.target.value)} />
          </div>
          <div className="space-y-2">
            <Label htmlFor="project-description">Descrição</Label>
            <Input id="project-description" value={description} onChange={(e) => setDescription(e.target.value)} />
          </div>
          <Button className="w-full" onClick={handleCreate}>Salvar</Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
