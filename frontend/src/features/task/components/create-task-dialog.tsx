import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

export function CreateTaskDialog({ onCreate, disabled }: { onCreate: (title: string) => Promise<void>; disabled?: boolean }) {
  const [title, setTitle] = useState('');

  const handleCreate = async () => {
    await onCreate(title);
    setTitle('');
  };

  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button disabled={disabled}>Nova tarefa</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Criar tarefa</DialogTitle>
          <DialogDescription>Adicione um item ao fluxo atual do projeto.</DialogDescription>
        </DialogHeader>
        <div className="space-y-4 pt-4">
          <div className="space-y-2">
            <Label htmlFor="task-title">Título</Label>
            <Input id="task-title" value={title} onChange={(event) => setTitle(event.target.value)} placeholder="Ex.: Ajustar fluxo de aprovação" />
          </div>
          <Button className="w-full" onClick={handleCreate} disabled={!title.trim()}>Salvar</Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
